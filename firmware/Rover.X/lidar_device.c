#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

#define RETRY_MAX  100  // define the retry count
#define LIDAR_ADDRESS    0x62 // slave device address
#define MAG_ADDRESS (0x1E)

extern Queue CmdQ[MAX_DEVICES];
void Mag_init(uint8_t reg, uint8_t byte);
uint8_t Mag_Read(uint16_t address, uint8_t *pData, uint16_t nCount);

uint8_t Lidar_Read(
        uint16_t address,
        uint8_t *pData,
        uint16_t nCount);

void LidarInit(void) {
  uint8_t data[2] = {0, 0};
  I2C1_MESSAGE_STATUS status;

  data[0] = 0x02;
  data[1] = 0x80;
  I2C1_MasterWrite(data, 2, LIDAR_ADDRESS, &status);
  while (status != I2C1_MESSAGE_COMPLETE);
   __delay_ms(10);
   
  data[0] = 0x04;
  data[1] = 0x08;
  I2C1_MasterWrite(data, 2, LIDAR_ADDRESS, &status);
  while (status != I2C1_MESSAGE_COMPLETE);
   __delay_ms(10);
  

  data[0] = 0x1C;
  data[1] = 0xB0;
  I2C1_MasterWrite(data, 2, LIDAR_ADDRESS, &status);
  while (status != I2C1_MESSAGE_COMPLETE);
   __delay_ms(10);
}

void LidarTask(void) {
  if (CmdQ[DEV_LIDAR].free) {
    return;
  }
  uint8_t data[2] = {0, 0};
  I2C1_MESSAGE_STATUS status;

  uint8_t command, packet[PKT_SZ];
  command = GetCommand(CmdQ[DEV_LIDAR].packet[0]);

  switch (command) {

    case CMD_STATE:
      //Lidar_Read(0x16, data, 1);
      // Mag_init(0x00, 0x04);

      // Start Acq
      /*  data[0] = 0x00;
        data[1] = 0x00;
        I2C1_MasterWrite(data, 2, LIDAR_ADDRESS, &status);
        while (status != I2C1_MESSAGE_COMPLETE);*/

      // Start Acq
      data[0] = 0x00;
      data[1] = 0x03;
      I2C1_MasterWrite(data, 2, LIDAR_ADDRESS, &status);
      while (status != I2C1_MESSAGE_COMPLETE);
      // __delay_ms(10);

      /*   // Read status reg.
         data[0] = 0x01;
         I2C1_MasterWrite(data, 1, LIDAR_ADDRESS, &status);
         while (status != I2C1_MESSAGE_COMPLETE);
        // __delay_ms(10);
         I2C1_MasterRead(data, 1, LIDAR_ADDRESS, &status);
         while (status != I2C1_MESSAGE_COMPLETE);*/
      __delay_ms(50);

      // Read value reg.
      data[0] = 0x8f;
      I2C1_MasterWrite(data, 1, LIDAR_ADDRESS, &status);
      while (status != I2C1_MESSAGE_COMPLETE);
      // __delay_ms(10);
      I2C1_MasterRead(data, 2, LIDAR_ADDRESS, &status);
      while (status != I2C1_MESSAGE_COMPLETE);
      // __delay_ms(10);

      // Read value reg.
      data[0] = 0x8f;
      I2C1_MasterWrite(data, 1, LIDAR_ADDRESS, &status);
      while (status != I2C1_MESSAGE_COMPLETE);
      // __delay_ms(10);
      I2C1_MasterRead(data, 2, LIDAR_ADDRESS, &status);
      while (status != I2C1_MESSAGE_COMPLETE);
      // __delay_ms(10);


      LED1_SetHigh();

      //  Mag_Read(0x7,data,2);
      packet[0] = 0xC0; // Ack & Done.
      packet[1] = data[0];
      packet[2] = data[1];
      SendPacket(DEV_LIDAR, packet, 3);
      break;

    default:
      SendError(DEV_LIDAR, ERR_UNIMPLEMENTED);
      break;

  }
  CmdQ[DEV_LIDAR].free = true;
}

void Mag_init(uint8_t reg, uint8_t byte) {
  I2C1_MESSAGE_STATUS status;
  uint8_t writeBuffer[3];
  uint16_t timeOut = 0;

  // writeBuffer[0] = (uint8_t) (reg >> 8); // high address
  writeBuffer[0] = reg; // low low address
  writeBuffer[1] = byte; // low low address

  while (status != I2C1_MESSAGE_FAIL) {
    I2C1_MasterWrite(writeBuffer, 2, MAG_ADDRESS, &status);
    // wait for the message to be sent or status has changed.
    while (status == I2C1_MESSAGE_PENDING);

    if (status == I2C1_MESSAGE_COMPLETE) {
      return;
    }

    // if status is  I2C1_MESSAGE_ADDRESS_NO_ACK,
    //               or I2C1_DATA_NO_ACK,
    // The device may be busy and needs more time for the last
    // write so we can retry writing the data, this is why we
    // use a while loop here

    // check for max retry and skip this byte
    if (timeOut == RETRY_MAX) {

      return;
    } else
      timeOut++;
  }

}

uint8_t Mag_Read(uint16_t address, uint8_t *pData, uint16_t nCount) {
  I2C1_MESSAGE_STATUS status;
  uint8_t writeBuffer[6];
  uint16_t timeOut;
  uint16_t counter;
  uint8_t *pD, ret;

  pD = pData;

  for (counter = 0; counter < nCount; counter++) {

    // build the write buffer first
    // starting address of the EEPROM memory
    writeBuffer[0] = (address); // high address
    // writeBuffer[1] = (uint8_t) (address); // low low address

    // Now it is possible that the slave device will be slow.
    // As a work around on these slaves, the application can
    // retry sending the transaction
    timeOut = 0;
    while (status != I2C1_MESSAGE_FAIL) {
      // write one byte to EEPROM (2 is the count of bytes to write)
      I2C1_MasterWrite(writeBuffer,
              1,
              MAG_ADDRESS,
              &status);

      // wait for the message to be sent or status has changed.
      while (status == I2C1_MESSAGE_PENDING);

      if (status == I2C1_MESSAGE_COMPLETE)
        break;

      // if status is  I2C1_MESSAGE_ADDRESS_NO_ACK,
      //               or I2C1_DATA_NO_ACK,
      // The device may be busy and needs more time for the last
      // write so we can retry writing the data, this is why we
      // use a while loop here

      // check for max retry and skip this byte
      if (timeOut == 100)
        break;
      else
        timeOut++;
    }

    if (status == I2C1_MESSAGE_COMPLETE) {

      // this portion will read the byte from the memory location.
      timeOut = 0;
      while (status != I2C1_MESSAGE_FAIL) {
        // write one byte to EEPROM (2 is the count of bytes to write)
        I2C1_MasterRead(pD,
                1,
                MAG_ADDRESS,
                &status);

        // wait for the message to be sent or status has changed.
        while (status == I2C1_MESSAGE_PENDING);

        if (status == I2C1_MESSAGE_COMPLETE) {
          break;
        }
        // if status is  I2C1_MESSAGE_ADDRESS_NO_ACK,
        //               or I2C1_DATA_NO_ACK,
        // The device may be busy and needs more time for the last
        // write so we can retry writing the data, this is why we
        // use a while loop here

        // check for max retry and skip this byte
        if (timeOut == 100) {
          break;
        } else
          timeOut++;
      }
    }

    // exit if the last transaction failed
    if (status == I2C1_MESSAGE_FAIL) {
      ret = 0;
      break;
    }

    pD++;
    address++;

  }
  return (ret);

}

uint8_t Lidar_Read(uint16_t address, uint8_t *pData, uint16_t nCount) {
  I2C1_MESSAGE_STATUS status;
  uint8_t writeBuffer[3];
  uint16_t timeOut;
  uint16_t counter;
  uint8_t *pD, ret;

  pD = pData;

  for (counter = 0; counter < nCount; counter++) {

    // build the write buffer first
    // starting address of the EEPROM memory
    writeBuffer[0] = (address); // high address
    //   writeBuffer[1] = (uint8_t) (address); // low low address
    // Now it is possible that the slave device will be slow.
    // As a work around on these slaves, the application can
    // retry sending the transaction
    timeOut = 0;
    while (status != I2C1_MESSAGE_FAIL) {
      // write one byte to EEPROM (2 is the count of bytes to write)
      I2C1_MasterWrite(writeBuffer,
              1,
              LIDAR_ADDRESS,
              &status);
      // wait for the message to be sent or status has changed.
      while (status == I2C1_MESSAGE_PENDING);
      LED1_SetHigh();

      if (status == I2C1_MESSAGE_COMPLETE)
        break;

      // if status is  I2C1_MESSAGE_ADDRESS_NO_ACK,
      //               or I2C1_DATA_NO_ACK,
      // The device may be busy and needs more time for the last
      // write so we can retry writing the data, this is why we
      // use a while loop here

      // check for max retry and skip this byte
      if (timeOut == RETRY_MAX)
        break;
      else
        timeOut++;
    }

    if (status == I2C1_MESSAGE_COMPLETE) {

      // this portion will read the byte from the memory location.
      timeOut = 0;
      while (status != I2C1_MESSAGE_FAIL) {
        // write one byte to EEPROM (2 is the count of bytes to write)
        I2C1_MasterRead(pD,
                1,
                LIDAR_ADDRESS,
                &status);

        // wait for the message to be sent or status has changed.
        while (status == I2C1_MESSAGE_PENDING);

        if (status == I2C1_MESSAGE_COMPLETE)
          break;

        // if status is  I2C1_MESSAGE_ADDRESS_NO_ACK,
        //               or I2C1_DATA_NO_ACK,
        // The device may be busy and needs more time for the last
        // write so we can retry writing the data, this is why we
        // use a while loop here

        // check for max retry and skip this byte
        if (timeOut == RETRY_MAX)
          break;
        else
          timeOut++;
      }
    }

    // exit if the last transaction failed
    if (status == I2C1_MESSAGE_FAIL) {
      ret = 0;
      break;
    }

    pD++;
    address++;

  }
  return (ret);

}

