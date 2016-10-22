#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

uint16_t se_m1_count, se_m2_count = 0;

void MotorTask(void) {
    uint16_t slots = 0;
    uint8_t packet[PKT_SZ];

    if (CmdQ[DEV_MOTOR].free) {
        // nothing to do
        return;
    }

    uint8_t command = GetCommand(CmdQ[DEV_MOTOR].packet[0]);

    switch (command) {
        case CMD_FWD:
            if (CmdQ[DEV_MOTOR].size != 3) {
                // Send insufficient param error 
                SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
                break;
            }
            se_m1_count, se_m2_count = 0;
            slots = CmdQ[DEV_MOTOR].packet[1];
            slots = slots << 8 | CmdQ[DEV_MOTOR].packet[2];
            MOTOR1_FWD_SetHigh();
            MOTOR1_BWD_SetLow();
            MOTOR2_FWD_SetHigh();
            MOTOR2_BWD_SetLow();
            SendAckDone(DEV_MOTOR);
            break;

        case CMD_BWD:
            if (CmdQ[DEV_MOTOR].size != 3) {
                // Send insufficient param error 
                SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
                break;
            }
            se_m1_count, se_m2_count = 0;
            slots = CmdQ[DEV_MOTOR].packet[1];
            slots = slots << 8 | CmdQ[DEV_MOTOR].packet[2];
            MOTOR1_BWD_SetHigh();
            MOTOR1_FWD_SetLow();
            MOTOR2_BWD_SetHigh();
            MOTOR2_FWD_SetLow();
            SendAckDone(DEV_MOTOR);
            break;

        case CMD_OFF:
            MOTOR1_BWD_SetLow();
            MOTOR1_FWD_SetLow();
            MOTOR2_FWD_SetLow();
            MOTOR2_BWD_SetLow();
            SendAckDone(DEV_MOTOR);
            break;

        case CMD_STATE:
            packet[0] = 0xC0 | DEV_MOTOR;
            packet[1] = se_m2_count >> 8;
            packet[2] = se_m2_count & 0xFF;
            SendPacket(packet, 3);
            break;

        default:
            SendError(DEV_MOTOR, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_MOTOR].free = true;
}


