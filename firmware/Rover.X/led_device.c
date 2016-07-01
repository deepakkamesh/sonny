#include <stdbool.h>
#include <stdlib.h>
#include "admin_device.h"
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void LedTask(void) {
    static uint32_t last;
    static bool blink;
    static bool countBlink;
    static uint16_t duration; // Default blink duration in milli secs.
    static uint8_t count; // Default number of times to blink.

    // Do some regular tasks.
    if (blink) {
        uint32_t now = TickGet();

        if ((now - last) / (TICK_MILLISECOND) >= duration) {
            LED1_Toggle();
            last = now;
            if (countBlink) {
                if (count <= 1) {
                    blink = false;
                } else {
                    count--;
                }
            }

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
            blink = false;
            LED1_SetLow();
            SendAckDone(DEV_LED1);
            break;

        case CMD_BLINK:
            LED1_SetLow();
            blink = true;
            countBlink = false;
            count = 0;
            last = 0;
            duration = 1000; //default duration.
            if (CmdQ[DEV_LED1].size >= 3) {
                duration = CmdQ[DEV_LED1].packet[1];
                duration = duration << 8 | CmdQ[DEV_LED1].packet[2];
            }
            if (CmdQ[DEV_LED1].size >= 4 && CmdQ[DEV_LED1].packet[3] > 0) {
                count = (CmdQ[DEV_LED1].packet[3])*2; // Multiply by 2 to account for on and off state.
                countBlink = true;
            }
            SendAckDone(DEV_LED1);
            break;

        default:
            SendError(DEV_LED1, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_LED1].free = true;
}

