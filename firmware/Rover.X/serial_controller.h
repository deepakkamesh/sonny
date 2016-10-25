/* 
 * File:   serial_controller.h
 * Author: dkg
 *
 * Created on June 14, 2016, 8:04 PM
 */
#ifndef SERIAL_CONTROLLER_H
#define	SERIAL_CONTROLLER_H

#include <stdint.h>
#include <stdbool.h>

#ifdef	__cplusplus
extern "C" {
#endif


void SerialReadTask(void);
bool VerifyCheckSum(uint8_t a[],uint8_t len, uint8_t chksum);
uint8_t CalcCheckSum(uint8_t a[], uint8_t len);
void SendError(uint8_t devID, uint8_t error);
void SendAck(uint8_t devID);
void SendPacket(uint8_t packet[], uint8_t size);
void SendAckDone(uint8_t devID);
void SendDone(uint8_t devID);


#ifdef	__cplusplus
}
#endif

#endif	/* SERIAL_CONTROLLER_H */

