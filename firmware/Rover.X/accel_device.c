#include <stdbool.h>
#include <stdlib.h>
#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void AccelTask(void) {
  if (CmdQ[DEV_ACCEL].free) {
    // nothing to do
    return;
  }
  uint8_t command, packet[PKT_SZ];
  command = GetCommand(CmdQ[DEV_ACCEL].packet[0]);

  uint16_t gX, gY, gZ = 0;

  switch (command) {

    case CMD_STATE:
      gX = ADC_GetConversion(AX);
      gY = ADC_GetConversion(AY);
      gZ = ADC_GetConversion(AZ);
      packet[0] = 0xC0; // Ack & Done.
      packet[1] = gX >> 8; // Pack 10 bit ADC Value.
      packet[2] = gX & 0xFF;
      packet[3] = gY >> 8;
      packet[4] = gY & 0xFF;
      packet[5] = gZ >> 8;
      packet[6] = gZ & 0xFF;
      SendPacket(DEV_ACCEL, packet, 7);
      break;

    default:
      SendError(DEV_ACCEL, ERR_UNIMPLEMENTED);
      break;

  }
  CmdQ[DEV_ACCEL].free = true;

}
