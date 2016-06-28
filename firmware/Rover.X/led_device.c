#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void LedTask(void) {
    static uint32_t last;
    static bool blink = false;
    static uint16_t duration; //default blink duration in milli secs.

    // Do some regular tasks.
    if (blink) {
        uint32_t now = 0;
        now = TickGet();
        if ((now - last)/(TICK_MILLISECOND )>= duration) {
            LED1_Toggle();
            last = now;
        }
    }


    if (CmdQ[DEV_LED1].free) {
        // nothing to do
        return;
    }
    uint8_t command;

    command = GetCommand(CmdQ[DEV_LED1].packet[0]);

    switch (command) {
        case CMD_ON:
            blink = false;
            LED1_SetHigh();
            SendAckDone(DEV_LED1);
            break;

        case CMD_OFF:
            LED1_SetLow();
            blink = false;
            SendAckDone(DEV_LED1);
            break;

        case CMD_BLINK:
            blink = true;
            last = 0;
            duration = 1000; //default duration.
            // Load duration if specified.
            if (CmdQ[DEV_LED1].size == 3) {
                duration = CmdQ[DEV_LED1].packet[1];
                duration = duration << 8 | CmdQ[DEV_LED1].packet[2];
            }
            SendAckDone(DEV_LED1);
            break;

        default:
            SendError(DEV_LED1, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_LED1].free = true;
}

