#include <xc.h>
#include <stdbool.h>
#include <stdint.h>
#include <string.h>
#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h" 


Queue CmdQ[MAX_DEVICES];
Queue SendQ[MAX_DEVICES];

// Local function declarations.
void I2C2_Callback(I2C2_SLAVE_DRIVER_STATUS i2c_bus_state);

typedef enum {
  SLAVE_NORMAL_DATA,
  SLAVE_DATA_ADDRESS,
} SLAVE_WRITE_DATA_TYPE;

void HostControllerInit(void) {
  I2C2_SetCallback(I2C2_Callback);
  // Init Command Queue
  uint8_t i;
  for (i = 0; i < MAX_DEVICES; i++) {
    CmdQ[i].free = true;
    SendQ[i].free = true;
    CmdQ[i].size = 0;
    SendQ[i].size = 0;
  }
}
bool headSnd = true;

void I2C2_Callback(I2C2_SLAVE_DRIVER_STATUS i2c_bus_state) {

  static uint8_t deviceID = 0, ptr = 0, PktSz = 0, chksum = 0;
  static uint8_t slaveWriteType = SLAVE_NORMAL_DATA;
  static bool headNow = false;
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
          headNow = true;
          ptr = 0;
          break;

        case SLAVE_NORMAL_DATA:
        default:
          if (headNow) {
            headNow = false;
            data = I2C2_slaveWriteData;
            PktSz = data >> 4;
            chksum = data & 0xF;
            break;
          }
          if (!CmdQ[deviceID].free) {
            SendError(deviceID, ERR_DEVICE_BUSY);
            break;
          }
          CmdQ[deviceID].packet[ptr++] = I2C2_slaveWriteData;
          if (PktSz == ptr) {
            if (!VerifyCheckSum(CmdQ[deviceID].packet, ptr, chksum)) {
              SendError(deviceID, ERR_CHECKSUM_FAILURE);
              break;
            }
            CmdQ[deviceID].size = ptr;
            CmdQ[deviceID].free = false;
          }
          break;
      } // end switch(slaveWriteType)

      slaveWriteType = SLAVE_NORMAL_DATA;
      break;

    case I2C2_SLAVE_READ_REQUEST:
      PktSz = SendQ[deviceID].size;
      // If free, nothing to send.
      if (SendQ[deviceID].free || ptr >= PktSz) {
        SSP2BUF = 0;
        break;
      }
      if (headSnd) {
        headSnd = false;
        chksum = CalcCheckSum(SendQ[deviceID].packet, PktSz); // 4 bit checksum.
        uint8_t header = PktSz << 4 | (chksum & 0xF);
        SSP2BUF = header;
        break;
      }
      SSP2BUF = SendQ[deviceID].packet[ptr++];
      if (PktSz == ptr) {
        SendQ[deviceID].free = true;
        headSnd = true;
      }
      break;

    case I2C2_SLAVE_READ_COMPLETED:
    default:;

  } // end switch(i2c_bus_state)
}

void SendError(uint8_t devID, uint8_t error) {
  uint8_t packet[PKT_SZ]; //minus header.
  packet[0] = 0;
  packet[1] = error;
  SendPacket(devID, packet, 2);
}

void SendAck(uint8_t devID) {
  uint8_t packet[PKT_SZ];
  packet[0] = 0x80;
  SendPacket(devID, packet, 1);
}

void SendAckDone(uint8_t devID) {
  uint8_t packet[PKT_SZ];
  packet[0] = 0xC0;
  SendPacket(devID, packet, 1);
}

void SendDone(uint8_t devID) {
  uint8_t packet[PKT_SZ];
  packet[0] = 0x40;
  SendPacket(devID, packet, 1);
}

void SendPacket(uint8_t deviceID, uint8_t packet[], uint8_t size) {

  headSnd = true;
  SendQ[deviceID].size = size;
  memcpy(SendQ[deviceID].packet, packet, size);
  SendQ[deviceID].free = false;

}

uint8_t CalcCheckSum(uint8_t a[], uint8_t len) {
  int sum = 0;
  for (uint8_t i = 0; i < len; i++) {
    sum += a[i]*(i + 1);
  }
  return sum % 16;
}

bool VerifyCheckSum(uint8_t a[], uint8_t len, uint8_t chksum) {
  int sum = 0;
  for (uint8_t i = 0; i < len; i++) {
    sum += a[i]*(i + 1);
  }

  return ((sum % 16) == chksum);
}
