#include <stdbool.h>
#include <stdlib.h>
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void EdgeSensorTask(void) {
    uint8_t state=0;
    uint32_t now;
    static uint32_t timer = 0;



    if (CmdQ[DEV_EDGE_SENSOR].free) {
        // nothing to do
        return;
    }
    uint8_t command, packet[PKT_SZ];
    command = GetCommand(CmdQ[DEV_EDGE_SENSOR].packet[0]);


    switch (command) {

        default:
            SendError(DEV_EDGE_SENSOR, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_EDGE_SENSOR].free = true;
}

