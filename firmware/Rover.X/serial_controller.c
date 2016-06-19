#include <xc.h>
#include <stdbool.h>
#include <stdint.h>
#include <string.h>
#include "serial_controller.h"
#include "protocol.h"
#include "mcc_generated_files/mcc.h" 
// SerialTask runs in a loop to process incoming serial data and 
// hands it over to the respective device to handle.
extern Queue CmdQ[MAX_DEVICES];

void SerialReadTask(void) {

    static uint8_t data, packetSize, chksum, sz, packet[PKT_SZ], deviceID;


    // Check if there is data to be read.
    RCSTA1bits.SREN = 1;
    if (!PIR1bits.RC1IF) {
        return;
    }
    if (1 == RCSTA1bits.OERR) {
        // EUSART1 error - restart
        RCSTA1bits.SPEN = 0;
        RCSTA1bits.SPEN = 1;
    }
    data = RCREG1;

    // State machine definitions.

    static enum {
        NEW_PACKET = 0,
        ASSEMBLE_PACKET,
        PROCESS_PACKET,
    } state = NEW_PACKET;

    switch (state) {

        case NEW_PACKET:
            sz = 0;
            packetSize = data >> 4;
            chksum = data & 0xF;
            state = ASSEMBLE_PACKET;
            break;

        case ASSEMBLE_PACKET:
            packet[sz] = data;
            sz++;
            if (sz == packetSize) {
                if (!VerifyCheckSum(packet, sz, chksum)) {
                    // Checksum failed. Discard packet and send error.
                    SendError(0, ERR_CHECKSUM_FAILURE);
                    state = NEW_PACKET;
                    break;
                }
                state = PROCESS_PACKET;
            } else {
                break;
            }

        case PROCESS_PACKET:
            deviceID = packet[0] & 0xF;
            // Add command to queue if device is free.
            if (CmdQ[deviceID].free) {
                memcpy(CmdQ[deviceID].packet, packet, packetSize);
                CmdQ[deviceID].size = packetSize;
                CmdQ[deviceID].free = false;
            } else {
                SendError(deviceID, ERR_DEVICE_BUSY);
            }
            state = NEW_PACKET;
            break;
        default:
            break;
    }
}

void SendError(uint8_t devID, uint8_t error) {
    uint8_t packet[PKT_SZ]; //minus header.
    packet[0] = (devID & 0xF);
    packet[1] = error;
    SendPacket(packet, 2);
}

void SendAck(uint8_t devID) {
    uint8_t packet[PKT_SZ];
    packet[0] = 0x80 | devID;
    SendPacket(packet, 1);
}

void SendAckDone(uint8_t devID) {
    uint8_t packet[PKT_SZ];
    packet[0] = 0xC0 | devID;
    SendPacket(packet, 1);
}

void SendPacket(uint8_t packet[], uint8_t size) {

    // Calculate checksum of packet.
    uint8_t chksum, header;
    chksum = CalcCheckSum(packet, size); // 4 bit checksum.
    header = size << 4 | (chksum & 0xF);

    EUSART1_Write(header);
    uint8_t i;
    for (i = 0; i < size; i++) {
        EUSART1_Write(packet[i]);
    }
}

uint8_t CalcCheckSum(uint8_t a[], uint8_t len) {
    return 0x6;
}

bool VerifyCheckSum(uint8_t a[], uint8_t len, uint8_t chksum) {
    return true;
}
