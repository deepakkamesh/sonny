#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "tick.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];
unsigned short t2OF = 0; // Timer 2 Overflow.

// myTMR2ISR is the callback function for Timer2 Interrupt.

void myTMR2ISR(void) {
  t2OF = 1;
}

// DHT11Init handles any initialization for DHT11.

void DHT11Init(void) {
  TMR2_SetInterruptHandler(myTMR2ISR);
}

// ReadByte reads 1 byte from DHT11 sensor. It returns -1 on failure.

int16_t ReadByte() {
  uint8_t data, dur;

  for (uint8_t i = 0; i < 8; i++) {
    t2OF = 0;
    TMR2_WriteTimer(0);
    while (!DHT11_GetValue() && !t2OF); // 50us low preamble.
    if (t2OF) {
      return -1;
    }

    TMR2_WriteTimer(0);
    while (DHT11_GetValue() && !t2OF); // Read bit type.
    if (t2OF) {
      return -1;
    }
    dur = TMR2_ReadTimer();
    if (dur < 40) {
      data = data << 1;
    } else {
      data = data << 1 | 0x1;
    }
  }
  return data;
}

// DHT11Task handles requests for Temp and Humidity.

void DHT11Task(void) {
  uint32_t now, delta = 0;
  uint8_t command, packet[PKT_SZ];

  if (CmdQ[DEV_DHT11].free) {
    // nothing to do
    return;
  }

  command = GetCommand(CmdQ[DEV_DHT11].packet[0]);

  switch (command) {
    case CMD_STATE:

      // Send start signal.
      DHT11_SetDigitalOutput();
      DHT11_SetLow();
      __delay_ms(25);
      DHT11_SetHigh();
      __delay_us(30);
      DHT11_SetDigitalInput();

      // Check response sequence from sensor.
      TMR2_WriteTimer(0);
      t2OF = 0;
      while (!DHT11_GetValue() && !t2OF); // 80us low.
      if (t2OF) {
        SendError(DEV_DHT11, ERR_TIMEOUT);
        break;
      }
      TMR2_WriteTimer(0);
      while (DHT11_GetValue() && !t2OF); // 80us high.
      if (t2OF) {
        SendError(DEV_DHT11, ERR_TIMEOUT);
        break;
      }
      
      // Read data from sensor.
      for (uint8_t i = 1; i <= 5; i++) {
        int16_t data = ReadByte();
        if (data == -1) {
          SendError(DEV_DHT11, ERR_TIMEOUT);
          break;
        }
        packet[i] = data;
      }
      
      // Verify checksum.
      if (packet[5] != (packet[1] + packet[2] + packet[3] + packet[4])& 0xFF) {
        SendError(DEV_DHT11, ERR_CHECKSUM_FAILURE);
        break;
      }
      
      packet[0] = 0xC0 | DEV_DHT11; // Ack & Done.
      SendPacket(packet, 5);
      break;

    default:
      SendError(DEV_DHT11, ERR_UNIMPLEMENTED);
      break;
  }

  CmdQ[DEV_DHT11].free = true;
}

