#include <stdbool.h>
#include <stdlib.h>
#include "host_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h"

extern Queue CmdQ[MAX_DEVICES];

void LDRTask(void) {


    if (CmdQ[DEV_LDR].free) {
        // nothing to do
        return;
    }
    uint8_t command, packet[PKT_SZ];
    command = GetCommand(CmdQ[DEV_LDR].packet[0]);
    uint16_t ldr = 0;

    switch (command) {
        case CMD_STATE:
            ldr = ADC_GetConversion(LDR);
            packet[0] = 0xC0 | DEV_LDR; // Ack & Done.
            packet[1] = ldr >> 8; // Pack 10 bit ADC Value.
            packet[2] = ldr & 0xFF;
            SendPacket(packet, 3);
            break;

        default:
            SendError(DEV_LDR, ERR_UNIMPLEMENTED);
            break;
    }

    CmdQ[DEV_LDR].free = true;
}


