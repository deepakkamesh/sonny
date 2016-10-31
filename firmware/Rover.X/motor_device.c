#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

uint16_t se_m1_count = 0, se_m2_count = 0;

void MotorTask(void) {
    uint8_t packet[PKT_SZ];
    static uint16_t rotation = 0;
    static bool active = false;

    // Check if rotations are done.
    if (se_m1_count / 40 >= rotation && active) {
        MOTOR1_BWD_SetLow();
        MOTOR1_FWD_SetLow();
        MOTOR2_FWD_SetLow();
        MOTOR2_BWD_SetLow();
        SendDone(DEV_MOTOR);
        active = false;

    }
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
            // Initialize counters.
            se_m1_count = 0;
            se_m2_count = 0;
            rotation = CmdQ[DEV_MOTOR].packet[1];
            rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
            MOTOR1_FWD_SetHigh();
            MOTOR1_BWD_SetLow();
            MOTOR2_FWD_SetHigh();
            MOTOR2_BWD_SetLow();
            SendAck(DEV_MOTOR);
            active= true;
            break;

        case CMD_BWD:
            if (CmdQ[DEV_MOTOR].size != 3) {
                // Send insufficient param error 
                SendError(DEV_MOTOR, ERR_INSUFFICENT_PARAMS);
                break;
            }
            se_m1_count = 0;
            se_m2_count = 0;
            rotation = CmdQ[DEV_MOTOR].packet[1];
            rotation = rotation << 8 | CmdQ[DEV_MOTOR].packet[2];
            MOTOR1_BWD_SetHigh();
            MOTOR1_FWD_SetLow();
            MOTOR2_BWD_SetHigh();
            MOTOR2_FWD_SetLow();
            SendAck(DEV_MOTOR);
            active = true;
            break;

        case CMD_OFF:
            MOTOR1_BWD_SetLow();
            MOTOR1_FWD_SetLow();
            MOTOR2_FWD_SetLow();
            MOTOR2_BWD_SetLow();
            SendAckDone(DEV_MOTOR);
            active = false;
            break;

        case CMD_STATE:
            packet[0] = 0xC0 | DEV_MOTOR;
            packet[1] = se_m1_count >> 8;
            packet[2] = se_m1_count & 0xFF;
            packet[3] = se_m2_count >> 8;
            packet[4] = se_m2_count & 0xFF;
            SendPacket(packet, 5);
            break;

        default:
            SendError(DEV_MOTOR, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_MOTOR].free = true;
}

// Called from IOC pin_manager.c

void SpeedEncoderISR_M1(void) {
    if (SE_M1_GetValue() == 0) {
        se_m1_count++;
        NOP();
    }
}

void SpeedEncoderISR_M2(void) {
    if (SE_M2_GetValue() == 0) {
        se_m2_count++;
        NOP();
    }
}