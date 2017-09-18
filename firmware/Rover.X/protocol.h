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
#define DEV_BATT 0x6
  /* Battery Voltage
   * Additional Parameters
   * CMD_STATE
   * Return:
   * Param1 - MSB of ADC Value 
   * Param2 - LSB of ADC Value
   */
#define DEV_MOTOR 0x7
  /* Drive Motor
   * Additional Parameters
   * CMD_FWD
   * Request:
   * Param1: MSB of slots to move
   * Param2: LSB of slots to move
   * 
   * CMD_BWD
   * Request:
   * Param1: MSB of slots to move
   * Param2: LSB of slots to move
   * 
   * CMD_STATE
   * Return:
   * Byte1: MSB of Motor1 slots moved.
   * Byte2: LSB of Motor1 slots moved.
   * Byte3: MSB of Motor2 slots moved.
   * Byte4: LSB of Motor2 slots moved.
   */
#define DEV_DHT11 0x8
  /* DHT11 Humidity/ Temp sensor
   * CMD_STATE
   * Return
   * Byte1 - MSB of Humidity
   * Byte2 - LSB of Humidity
   * Byte3 - MSB of Temp
   * Byte4 - LSB of Temp
   */
#define DEV_US020 0x9
  /* US020 Ultrasonic range finder.
   * CMD_STATE:
   * Return
   * Byte1 - MSB of distance in centimeters
   * Byte2 - LSB of distance in centimeters.
   */
#define DEV_LIDAR 0xA
  /* Lidar Garmin V3
   * CMD_STATE:
   * Return
   * 
   */
  // Command definitions.    
#define CMD_ON 0x1
#define CMD_PING 0x2
#define CMD_VERSION 0x3
#define CMD_OFF 0x4
#define CMD_BLINK 0x5 
#define CMD_ROTATE 0x6 
#define CMD_STATE 0x7
#define CMD_TEST 0x8
#define CMD_FWD 0x9
#define CMD_BWD 0xA

  // Error Codes.
#define ERR_CHECKSUM_FAILURE 0x1
#define ERR_DEVICE_BUSY 0x2    
#define ERR_UNIMPLEMENTED 0x3
#define ERR_INSUFFICENT_PARAMS 0x4
#define ERR_EDGE_DETECTED 0x5
#define ERR_BATT_LOW 0x6
#define ERR_TIMEOUT 0x7

  // Helper Functions.
#define GetDeviceID(data) (data & 0xF)
#define GetCommand(data) (data)


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

