#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void ServoTask(void) {
     if (CmdQ[DEV_SERVO].free) {
        // nothing to do
        return;
    }
    uint8_t command;
    command = GetCommand(CmdQ[DEV_SERVO].packet[0]);
    uint16_t on;
    switch (command) {

        case CMD_ROTATE:
            on = 1000; //default duration. Center.
            // Load on time. TODO set limits.
            if (CmdQ[DEV_SERVO].size == 3) {
                on = CmdQ[DEV_SERVO].packet[1];
                on = on << 8 | CmdQ[DEV_SERVO].packet[2];
            }
            CCP4_SetOnOff(on, 10000 - on);
            SendAckDone(DEV_SERVO);
            break;

        default:
            SendError(DEV_SERVO, ERR_UNIMPLEMENTED);
            break;

    }
    CmdQ[DEV_SERVO].free = true;
}