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
    /* Additional Parameters
     * CMD_BLINK
     * Request:
     * param1 - MSB of blink duration in ms 
     * param2 - LSB of blink duration in ms.
     * defaut - 1 sec
     */
#define DEV_SERVO 0x2
    /* Additional Parameters
     * CMD_ROTATE
     * Request:
     * param1 - MSB of pwn high duration in cycles.
     * param2 - LSB of pwn high duration in cycles.
     * param3 - Servo select 0x1 - Servo1 0x2 - Servo2
     * defaut - ?
     */
#define DEV_ACCEL 0x3
    /* Additional Parameters
     * CMD_STATE
     * Return:
     * Param1 - MSB of X Axis 
     * Param2 - LSB of X Axis
     * Param3 - MSB of Y Axis
     * Param4 - LSB of Y Axis
     * Param5 - MSB of Z Axis
     * Param6 - LSB of Z Axis
     */

    // Command definitions.    
#define CMD_ON 0x1
#define CMD_PING 0x2
#define CMD_VERSION 0x3
#define CMD_OFF 0x4
#define CMD_BLINK 0x5 
#define CMD_ROTATE 0x6 
#define CMD_STATE 0x7

    // Error Codes.
#define ERR_CHECKSUM_FAILURE 0x1
#define ERR_DEVICE_BUSY 0x2    
#define ERR_UNIMPLEMENTED 0x3
#define ERR_INSUFFICENT_PARAMS 0x4

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

