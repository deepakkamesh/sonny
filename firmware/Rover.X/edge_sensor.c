#include <stdbool.h>
#include <stdlib.h>
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void EdgeSensorTask(void) {
    uint8_t state, now;
    static uint8_t timer = 0;

    state = IR1_GetValue();
    

    if (!state) {
        now = GetTicks();
        if (!timer) {
            timer = now;
        } else {
            if (abs(now - timer) > 1) {
                SendError(DEV_EDGE_SENSOR, ERR_EDGE_DETECTED);
                timer = 0;
                return;
            }
        }
    } else {
        timer = 0;
    }



    if (CmdQ[DEV_EDGE_SENSOR].free) {
        // nothing to do
        return;
    }
    uint8_t command, packet[PKT_SZ];
    command = GetCommand(CmdQ[DEV_EDGE_SENSOR].packet[0]);


    switch (command) {
        case CMD_STATE:
            packet[0] = 0xC0 | DEV_EDGE_SENSOR; // Ack & Done.
            packet[1] = state;
            SendPacket(packet, 2);
            break;

        default:
            SendError(DEV_EDGE_SENSOR, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_EDGE_SENSOR].free = true;
}

