/* 
 * File:   protocol.h
 * Author: dkg
 *
 * Created on June 14, 2016, 8:27 PM
 */

#ifndef PROTOCOL_H
#define	PROTOCOL_H

#ifdef	__cplusplus
extern "C" {
#endif
#define PKT_SZ 16
#define MAX_DEVICES 16

    // Device Definitions.
#define DEV_ADMIN 0x0
#define DEV_LED1 0x1
#define DEV_SERVO 0x2 
#define DEV_ACCEL 0x3

    // Command definitions.    
#define CMD_ON 0x1
#define CMD_PING 0x2
#define CMD_VERSION 0x3
#define CMD_OFF 0x4
#define CMD_BLINK 0x5 // param1 - MSB, param2 - LSB of blink duration in ms.
// CMD_ROTATE.
// param1 - MSB, param2 - LSB of pwn high duration in cycles.
#define CMD_ROTATE 0x6 
#define CMD_STATE 0x7

    // Error Codes.
#define ERR_CHECKSUM_FAILURE 0x1
#define ERR_DEVICE_BUSY 0x2    
#define ERR_UNIMPLEMENTED 0x3    

    // Helper Functions.
#define GetDeviceID(data) (data & 0xF)
#define GetCommand(data) (data>>4 & 0xF)


    // Global device command queue.

    typedef struct {
        uint8_t packet[PKT_SZ];
        bool free;
        uint8_t size; // Size of packet.
    } Queue;

#ifdef	__cplusplus
}
#endif

#endif	/* PROTOCOL_H */

