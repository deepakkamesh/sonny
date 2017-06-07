
#include "admin_device.h"
#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void BattTask(void) {

  // TODO: Check battery level every 30s and report error.

  if (CmdQ[DEV_BATT].free) {
    // nothing to do
    return;
  }
  uint8_t deviceID, command, packet[PKT_SZ];
  uint16_t batt;

  command = GetCommand(CmdQ[DEV_BATT].packet[0]);

  switch (command) {
    case CMD_STATE:
      batt = ADC_GetConversion(channel_FVRBuf2);
      packet[0] = 0xC0; // Ack & Done.
      packet[1] = batt >> 8; // Pack 10 bit ADC Value.
      packet[2] = batt & 0xFF;
      SendPacket(DEV_BATT, packet, 3);
      break;
    default:
      SendError(DEV_BATT, ERR_UNIMPLEMENTED);
      break;
  }

  CmdQ[DEV_BATT].free = true;
}

