#include "admin_device.h"
#include "host_controller.h"
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
    if (dur < 20) {
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
  uint8_t command;
  uint8_t packet[PKT_SZ]; // TODO: Why is this working without being static?
  static uint32_t ticks = 0;

  static enum {
    SEND_TRIGGER = 0,
    WAIT_HIGH,
    READ_RESPONSE,
    READ_DATA,
    VERIFY_CHECKSUM,
    SEND_RESPONSE,
    RESET,
    READY,
  } state = READY;

  switch (state) {
    case SEND_TRIGGER:
      DHT11_SetDigitalOutput();
      DHT11_SetLow();
      ticks = TickGet();
      state = WAIT_HIGH;
      break;

    case WAIT_HIGH:
      if ((TickGet() - ticks) / TICK_MILLISECOND < 25) {
        break;
      }
      DHT11_SetHigh();
      __delay_us(25);
      DHT11_SetDigitalInput();
      state = READ_RESPONSE;

    case READ_RESPONSE:
      // Check response sequence from sensor.
      TMR2_WriteTimer(0);
      t2OF = 0;
      while (!DHT11_GetValue() && !t2OF); // 80us low.
      if (t2OF) {
        SendError(DEV_DHT11, ERR_TIMEOUT);
        state = RESET;
        break;
      }
      TMR2_WriteTimer(0);
      while (DHT11_GetValue() && !t2OF); // 80us high.
      if (t2OF) {
        SendError(DEV_DHT11, ERR_TIMEOUT);
        state = RESET;
        break;
      }
      state = READ_DATA;

    case READ_DATA:
      // Read data from sensor.
      for (uint8_t i = 1; i <= 5; i++) {
        int16_t data = ReadByte();
        if (data == -1) {
          SendError(DEV_DHT11, ERR_TIMEOUT);
          state = RESET;
          break;
        }
        packet[i] = data;
      }
      state = VERIFY_CHECKSUM;
      break;

    case VERIFY_CHECKSUM:
      // Verify checksum.
      if (packet[5] != (packet[1] + packet[2] + packet[3] + packet[4])& 0xFF) {
        SendError(DEV_DHT11, ERR_CHECKSUM_FAILURE);
        state = RESET;
        break;
      }
      state = SEND_RESPONSE;
      break;


    case SEND_RESPONSE:
      packet[0] = 0xC0; // Ack & Done.
      SendPacket(DEV_DHT11, packet, 5);
      state = RESET;
      break;

    case RESET:
      CmdQ[DEV_DHT11].free = true;
      state = READY;
      break;

    case READY:
      break;
  }

  if (CmdQ[DEV_DHT11].free || state != READY) {
    // nothing to do
    return;
  }

  command = GetCommand(CmdQ[DEV_DHT11].packet[0]);

  switch (command) {
    case CMD_STATE:
      state = SEND_TRIGGER;
      /*   // Send start signal.
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

         packet[0] = 0xC0; // Ack & Done.
         SendPacket(DEV_DHT11, packet, 5); */
      break;

    default:
      SendError(DEV_DHT11, ERR_UNIMPLEMENTED);
      break;
  }

}

