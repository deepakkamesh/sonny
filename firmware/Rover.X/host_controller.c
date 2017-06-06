#include <xc.h>
#include <stdbool.h>
#include <stdint.h>
#include <string.h>
#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h" 
// SerialTask runs in a loop to process incoming serial data and 
// hands it over to the respective device to handle.
extern Queue CmdQ[MAX_DEVICES];

typedef enum {
  SLAVE_NORMAL_DATA,
  SLAVE_DATA_ADDRESS,
} SLAVE_WRITE_DATA_TYPE;

void I2C2_Callback(I2C2_SLAVE_DRIVER_STATUS i2c_bus_state) {

  static uint8_t deviceID = 0, sz = 0, PktSz = 0;
  static uint8_t slaveWriteType = SLAVE_NORMAL_DATA;
  static bool gotHeader = false;
  uint8_t data;


  switch (i2c_bus_state) {
    case I2C2_SLAVE_WRITE_REQUEST:
      // the master will be sending the eeprom address next
      slaveWriteType = SLAVE_DATA_ADDRESS;
      break;

    case I2C2_SLAVE_WRITE_COMPLETED:

      switch (slaveWriteType) {
        case SLAVE_DATA_ADDRESS:
          deviceID = I2C2_slaveWriteData;
          break;

        case SLAVE_NORMAL_DATA:
        default:
          if (!gotHeader) {
            gotHeader = true;
            data = I2C2_slaveWriteData;
            PktSz = data >> 4;
            sz = 0;
            // TODO: Implement checksum verification.
            break;
          }
          CmdQ[deviceID].packet[sz++] = I2C2_slaveWriteData;
          if (PktSz == sz) {
            CmdQ[deviceID].size = sz;
            CmdQ[deviceID].free = false;
            gotHeader = false;
          }
          break;
      }
      slaveWriteType = SLAVE_NORMAL_DATA;
      break;

    case I2C2_SLAVE_READ_REQUEST:
      /*  SSP2BUF = EEPROM_Buffer[deviceID++];
        if (sizeof (EEPROM_Buffer) <= deviceID) {
          deviceID = 0; // wrap to start of eeprom page
        }*/
      break;

    case I2C2_SLAVE_READ_COMPLETED:
    default:;

  } // end switch(i2c_bus_state)
}

void HostControllerInit(void) {
  I2C2_SetCallback(I2C2_Callback);
}

void HostControllerTask(void) {

  static uint8_t data, packetSize, chksum, sz, packet[PKT_SZ], deviceID;


  // Check if there is data to be read.
  // TODO: Update to I2C
  //RCSTA1bits.SREN = 1;
  if (!PIR1bits.RC1IF) {
    return;
  }
  if (1 == RCSTA1bits.OERR) {
    // EUSART1 error - restart
    RCSTA1bits.SPEN = 0;
    RCSTA1bits.SPEN = 1;
  }
  data = RCREG1;

  // State machine definitions.

  static enum {
    NEW_PACKET = 0,
    ASSEMBLE_PACKET,
    PROCESS_PACKET,
  } state = NEW_PACKET;

  switch (state) {

    case NEW_PACKET:
      sz = 0;
      packetSize = data >> 4;
      if (packetSize == 0) { // Zero sized packet is not allowed.
        SendError(0, ERR_CHECKSUM_FAILURE);
        state = NEW_PACKET;
        break;
      }
      chksum = data & 0xF;
      state = ASSEMBLE_PACKET;
      break;

    case ASSEMBLE_PACKET:
      packet[sz] = data;
      sz++;
      if (sz == packetSize) {
        if (!VerifyCheckSum(packet, sz, chksum)) {
          // Checksum failed. Discard packet and send error.
          SendError(0, ERR_CHECKSUM_FAILURE);
          state = NEW_PACKET;
          break;
        }
        state = PROCESS_PACKET;
      } else {
        break;
      }

    case PROCESS_PACKET:
      deviceID = packet[0] & 0xF;
      // Add command to queue if device is free.
      if (CmdQ[deviceID].free) {
        memcpy(CmdQ[deviceID].packet, packet, packetSize);
        CmdQ[deviceID].size = packetSize;
        CmdQ[deviceID].free = false;
      } else {
        SendError(deviceID, ERR_DEVICE_BUSY);
      }
      state = NEW_PACKET;
      break;
    default:
      break;
  }
}

void SendError(uint8_t devID, uint8_t error) {
  uint8_t packet[PKT_SZ]; //minus header.
  packet[0] = (devID & 0xF);
  packet[1] = error;
  SendPacket(packet, 2);
}

void SendAck(uint8_t devID) {
  uint8_t packet[PKT_SZ];
  packet[0] = 0x80 | devID;
  SendPacket(packet, 1);
}

void SendAckDone(uint8_t devID) {
  uint8_t packet[PKT_SZ];
  packet[0] = 0xC0 | devID;
  SendPacket(packet, 1);
}

void SendDone(uint8_t devID) {
  uint8_t packet[PKT_SZ];
  packet[0] = 0x40 | devID;
  SendPacket(packet, 1);
}

void SendPacket(uint8_t packet[], uint8_t size) {

  // Calculate checksum of packet.
  uint8_t chksum, header;
  chksum = CalcCheckSum(packet, size); // 4 bit checksum.
  header = size << 4 | (chksum & 0xF);

  // TODO: Write Data to I2C
  // EUSART1_Write(header);
  uint8_t i;
  for (i = 0; i < size; i++) {
    //EUSART1_Write(packet[i]);
  }
}

uint8_t CalcCheckSum(uint8_t a[], uint8_t len) {
  return 0x6;
}

bool VerifyCheckSum(uint8_t a[], uint8_t len, uint8_t chksum) {
  return true;
}
