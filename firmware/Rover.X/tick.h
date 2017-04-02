/* 
 * File:   tick.h
 * Author: dkg
 *
 * Created on April 1, 2017, 10:17 PM
 */
#include <stdint.h>
#ifndef TICK_H
#define	TICK_H

#ifdef	__cplusplus
extern "C" {
#endif

void InitTicker(void);
void TMR0ISR(void);  
uint32_t TickGet(void);
static void GetTickCopy(void);
uint32_t TickGetDev64K(void);
uint32_t TickGetDev256(void);       

#define TICKS_PER_SECOND ((_XTAL_FREQ/4)/256ull)
#define TICK_MILLISECOND (TICKS_PER_SECOND/1000)
#define TICK_SECOND TICKS_PER_SECOND
#define TICK_MINUTE (TICKS_PER_SECOND*60ull)
#define TICK_HOUR (TICKS_PER_SECOND*3600ull)



#ifdef	__cplusplus
}
#endif

#endif	/* TICK_H */

