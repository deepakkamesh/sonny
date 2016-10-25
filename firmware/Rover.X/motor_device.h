/* 
 * File:   motor_device.h
 * Author: dkg
 *
 * Created on October 21, 2016, 8:15 PM
 */

#ifndef MOTOR_DEVICE_H
#define	MOTOR_DEVICE_H

#ifdef	__cplusplus
extern "C" {
#endif



void MotorTask(void);
void SpeedEncoderISR_M1(void);
void SpeedEncoderISR_M2(void);

#ifdef	__cplusplus
}
#endif

#endif	/* MOTOR_DEVICE_H */

