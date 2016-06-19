#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void AdminTask(void) {


    if (CmdQ[DEV_ADMIN].free) {
        // nothing to do
        return;
    }
    uint8_t deviceID, command, packet[PKT_SZ];
    unsigned long long now;

    command = GetCommand(CmdQ[DEV_ADMIN].packet[0]);

    switch (command) {
        case CMD_PING:
            SendAckDone(DEV_ADMIN);
            break;
            
        default:
            SendError(DEV_ADMIN, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_ADMIN].free = true;
}

