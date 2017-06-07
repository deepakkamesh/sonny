#include "admin_device.h"
#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void AdminTask(void) {


  if (CmdQ[DEV_ADMIN].free) {
    // nothing to do
    return;
  }
  uint8_t deviceID, command, packet[PKT_SZ];

  command = GetCommand(CmdQ[DEV_ADMIN].packet[0]);

  switch (command) {
    case CMD_PING:
      SendAckDone(DEV_ADMIN);
      break;
    case CMD_TEST:
      packet[0] = 0xC0; // Ack & Done.
      packet[1] = 0x40;
      packet[2] = 0x04;
      packet[3] = 0x30;
      packet[4] = 0x03;
      packet[5] = 0x20;
      packet[6] = 0x02;
      SendPacket(DEV_ADMIN,packet, 7);
      break;
    default:
      SendError(DEV_ADMIN, ERR_UNIMPLEMENTED);
      break;
  }

  CmdQ[DEV_ADMIN].free = true;
}

