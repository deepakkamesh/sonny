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
     * optional param1 - MSB of blink duration in ms 
     * optional param2 - LSB of blink duration in ms.
     * optional param3 - number of times to blink.
     * defaut - 1 sec blink continuous
     */
#define DEV_SERVO 0x2
    /* Servo Motors
     * Additional Parameters
     * CMD_ROTATE
     * Request:
     * param1 - MSB of pwn high duration in cycles.
     * param2 - LSB of pwn high duration in cycles.
     * param3 - MSB of pwn total period in cycles.
     * param4 - LSB of pwn total period in cycles.
     * param5 - Servo select 0x1 - Servo1 0x2 - Servo2
     * defaut - ?
     */
#define DEV_ACCEL 0x3
    /* Accelerometer
     * Additional Parameters
     * CMD_STATE
     * Return:
     * Param1 - MSB of X Axis 
     * Param2 - LSB of X Axis
     * Param3 - MSB of Y Axis
     * Param4 - LSB of Y Axis
     * Param5 - MSB of Z Axis
     * Param6 - LSB of Z Axis
     */
#define DEV_EDGE_SENSOR 0x4
    /* Edge Sensor
     * 
     */
#define DEV_LDR 0x5
    /* Light Sensor
     * Additional Parameters
     * CMD_STATE
     * Return:
     * Param1 - MSB of ADC Value 
     * Param2 - LSB of ADC Value
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
#define ERR_EDGE_DETECTED 0x5

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

