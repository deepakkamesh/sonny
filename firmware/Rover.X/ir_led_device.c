#include <stdbool.h>
#include <stdlib.h>
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void IRLedTask(void) {
    static uint32_t last;
    static bool blink = false;
    static uint16_t duration; //default blink duration in milli secs.

    // Do some regular tasks.
    if (blink) {
        uint32_t now = 0;
        now = TickGet();
        if ((now - last) / (TICK_MILLISECOND) >= duration) {
            PSTR1CON ^= 0x01; // Toggle PWM steering output.
            last = now;
        }
    }


    if (CmdQ[DEV_IR_LED].free) {
        // nothing to do
        return;
    }
    uint8_t command;

    command = GetCommand(CmdQ[DEV_IR_LED].packet[0]);

    switch (command) {
        case CMD_ON:
            blink = false;
            PSTR1CON = 0x01;
            SendAckDone(DEV_IR_LED);
            break;

        case CMD_OFF:
            PSTR1CON = 0x00;

            blink = false;
            SendAckDone(DEV_IR_LED);
            break;

        case CMD_BLINK:
            blink = true;
            last = 0;
            duration = 1000; //default duration.
            // Load duration if specified.
            if (CmdQ[DEV_IR_LED].size == 3) {
                duration = CmdQ[DEV_IR_LED].packet[1];
                duration = duration << 8 | CmdQ[DEV_IR_LED].packet[2];
            }
            SendAckDone(DEV_IR_LED);
            break;

        default:
            SendError(DEV_IR_LED, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_IR_LED].free = true;

}


