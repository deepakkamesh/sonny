/**
  TMR5 Generated Driver API Header File

  @Company
    Microchip Technology Inc.

  @File Name
    tmr5.h

  @Summary
    This is the generated header file for the TMR5 driver using MPLAB(c) Code Configurator

  @Description
    This header file provides APIs for driver for TMR5.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 4.15.3
        Device            :  PIC18F26K22
        Driver Version    :  2.00
    The generated drivers are tested against the following:
        Compiler          :  XC8 1.35
        MPLAB             :  MPLAB X 3.40
*/

/*
    (c) 2016 Microchip Technology Inc. and its subsidiaries. You may use this
    software and any derivatives exclusively with Microchip products.

    THIS SOFTWARE IS SUPPLIED BY MICROCHIP "AS IS". NO WARRANTIES, WHETHER
    EXPRESS, IMPLIED OR STATUTORY, APPLY TO THIS SOFTWARE, INCLUDING ANY IMPLIED
    WARRANTIES OF NON-INFRINGEMENT, MERCHANTABILITY, AND FITNESS FOR A
    PARTICULAR PURPOSE, OR ITS INTERACTION WITH MICROCHIP PRODUCTS, COMBINATION
    WITH ANY OTHER PRODUCTS, OR USE IN ANY APPLICATION.

    IN NO EVENT WILL MICROCHIP BE LIABLE FOR ANY INDIRECT, SPECIAL, PUNITIVE,
    INCIDENTAL OR CONSEQUENTIAL LOSS, DAMAGE, COST OR EXPENSE OF ANY KIND
    WHATSOEVER RELATED TO THE SOFTWARE, HOWEVER CAUSED, EVEN IF MICROCHIP HAS
    BEEN ADVISED OF THE POSSIBILITY OR THE DAMAGES ARE FORESEEABLE. TO THE
    FULLEST EXTENT ALLOWED BY LAW, MICROCHIP'S TOTAL LIABILITY ON ALL CLAIMS IN
    ANY WAY RELATED TO THIS SOFTWARE WILL NOT EXCEED THE AMOUNT OF FEES, IF ANY,
    THAT YOU HAVE PAID DIRECTLY TO MICROCHIP FOR THIS SOFTWARE.

    MICROCHIP PROVIDES THIS SOFTWARE CONDITIONALLY UPON YOUR ACCEPTANCE OF THESE
    TERMS.
*/

#ifndef _TMR5_H
#define _TMR5_H

/**
  Section: Included Files
*/

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus  // Provide C++ Compatibility

    extern "C" {

#endif


/**
  Section: TMR5 APIs
*/

/**
  @Summary
    Initializes the TMR5

  @Description
    This routine initializes the TMR5.
    This routine must be called before any other TMR5 routine is called.
    This routine should only be called once during system initialization.

  @Preconditions
    None

  @Param
    None

  @Returns
    None

  @Comment
    

  @Example
    <code>
    TMR5_Initialize();
    </code>
*/
void TMR5_Initialize(void);

/**
  @Summary
    Start TMR5

  @Description
    This routine is used to  Start TMR5.

  @Preconditions
    The TMR5_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    // Initialize TMR5 module

    // Start TMR5
    TMR5_StartTimer();

    while(1)
    {
    }
    </code>
 */
void TMR5_StartTimer(void);

/**
  @Summary
    Stop TMR5

  @Description
    This routine is used to  Stop TMR5.

  @Preconditions
    The TMR5_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    None

  @Example
    </code>
    TMR5_Initialize();

    TMR5_StartTimer();
    if(TMR5_HasOverflowOccured())
    {
        TMR5_StopTimer();
    }
    <code>
 */
void TMR5_StopTimer(void);

/**
  @Summary
    Read TMR5 register.

  @Description
    This routine is used to  Read TMR5 register.

  @Preconditions
    The TMR5_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    TMR5 value at the time of the function call read as a 16-bit value

  @Example
    <code>
    uint16_t timerVal=0;

    TMR5_Initialize();

    TMR5_StartTimer();

    // some delay or code

    TMR5_StopTimer();

    timerVal=TMR5_ReadTimer();
    </code>
 */
uint16_t TMR5_ReadTimer(void);

/**
  @Summary
    Write TMR5 register.

  @Description
    This routine is used to Write TMR5 register.

  @Preconditions
    The TMR5_Initialize() routine should be called
    prior to use this routine.

  @Param
    timerVal : Timer value to be loaded

  @Returns
    None

  @Example
    <code>
    TMR5_Initialize();
    TMR5_WriteTimer(0x055);
    TMR5_StartTimer();
    </code>
 */
void TMR5_WriteTimer(uint16_t timerVal);

/**
  @Summary
    Reload TMR5 register.

  @Description
    This routine is used to reload TMR5 register.

  @Preconditions
    The TMR5_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    TMR5_Initialize();
    TMR5_StartTimer();

    if(TMR5_HasOverflowOccured())
    {
        TMR5_StopTimer();
    }

    TMR5_Reload();}
    </code>
*/
void TMR5_Reload(void);

/**
  @Summary
    Starts the single pulse acquisition in TMR5 gate operation.

  @Description
    This function starts the single pulse acquisition in TMR5 gate operation.
    This function must be used when the TMR5 gate is enabled.

  @Preconditions
    Initialize  the TMR5 with gate enable before calling this function.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    uint16_t xVal;
    uint16_t yVal;

    //enable TMR5 singlepulse mode
    TMR5_StartSinglePulseAcquistion();

    //check TMR5 gate status
    if(TMR5_CheckGateValueStatus()== 0)
    {
        xVal = TMR5_ReadTimer();
    }

    // wait untill gate interrupt occured
    while(TMR5GIF == 0)
    {
    }

    yVal = TMR5_ReadTimer();
    </code>
*/
void TMR5_StartSinglePulseAcquisition(void);

/**
  @Summary
    Check the current state of Timer1 gate.

  @Description
    This function reads the TMR5 gate value and return it.
    This function must be used when the TMR5 gate is enabled.

  @Preconditions
    Initialize  the TMR5 with gate enable before calling this function.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    uint16_t xVal;
    uint16_t yVal;

    //enable TMR5 singlepulse mode
    TMR5_StartSinglePulseAcquistion();

    //check TMR5 gate status
    if(TMR5_CheckGateValueStatus()== 0)
    {
        xVal = TMR5_ReadTimer();
    }

    //wait untill gate interrupt occured
    while(TMR5IF == 0)
    {
    }

    yVal = TMR5_ReadTimer();
    </code>
*/
uint8_t TMR5_CheckGateValueStatus(void);

/**
  @Summary
    Get the TMR5 overflow status.

  @Description
    This routine get the TMR5 overflow status.

  @Preconditions
    The TMR5_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    true  - Overflow has occured.
    false - Overflow has not occured.

  @Comment


  @Example
     <code>
     TMR5_Initialize();

     TMR5_StartTimer();

    while(1)
    {
        if(TMR5_HasOverflowOccured())
        {
            TMR5_StopTimer();
        }
    }
    </code>
*/
bool TMR5_HasOverflowOccured(void);


#ifdef __cplusplus  // Provide C++ Compatibility

    }

#endif

#endif // _TMR5_H
/**
 End of File
*/