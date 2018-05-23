/**
  TMR3 Generated Driver API Header File

  @Company
    Microchip Technology Inc.

  @File Name
    tmr3.h

  @Summary
    This is the generated header file for the TMR3 driver using MPLAB(c) Code Configurator

  @Description
    This header file provides APIs for driver for TMR3.
    Generation Information :
        Product Revision  :  MPLAB(c) Code Configurator - 4.15.6
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

#ifndef _TMR3_H
#define _TMR3_H

/**
  Section: Included Files
*/

#include <stdbool.h>
#include <stdint.h>

#ifdef __cplusplus  // Provide C++ Compatibility

    extern "C" {

#endif


/**
  Section: TMR3 APIs
*/

/**
  @Summary
    Initializes the TMR3

  @Description
    This routine initializes the TMR3.
    This routine must be called before any other TMR3 routine is called.
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
    TMR3_Initialize();
    </code>
*/
void TMR3_Initialize(void);

/**
  @Summary
    Start TMR3

  @Description
    This routine is used to  Start TMR3.

  @Preconditions
    The TMR3_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    // Initialize TMR3 module

    // Start TMR3
    TMR3_StartTimer();

    while(1)
    {
    }
    </code>
 */
void TMR3_StartTimer(void);

/**
  @Summary
    Stop TMR3

  @Description
    This routine is used to  Stop TMR3.

  @Preconditions
    The TMR3_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    None

  @Example
    </code>
    TMR3_Initialize();

    TMR3_StartTimer();
    if(TMR3_HasOverflowOccured())
    {
        TMR3_StopTimer();
    }
    <code>
 */
void TMR3_StopTimer(void);

/**
  @Summary
    Read TMR3 register.

  @Description
    This routine is used to  Read TMR3 register.

  @Preconditions
    The TMR3_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    TMR3 value at the time of the function call read as a 16-bit value

  @Example
    <code>
    uint16_t timerVal=0;

    TMR3_Initialize();

    TMR3_StartTimer();

    // some delay or code

    TMR3_StopTimer();

    timerVal=TMR3_ReadTimer();
    </code>
 */
uint16_t TMR3_ReadTimer(void);

/**
  @Summary
    Write TMR3 register.

  @Description
    This routine is used to Write TMR3 register.

  @Preconditions
    The TMR3_Initialize() routine should be called
    prior to use this routine.

  @Param
    timerVal : Timer value to be loaded

  @Returns
    None

  @Example
    <code>
    TMR3_Initialize();
    TMR3_WriteTimer(0x055);
    TMR3_StartTimer();
    </code>
 */
void TMR3_WriteTimer(uint16_t timerVal);

/**
  @Summary
    Reload TMR3 register.

  @Description
    This routine is used to reload TMR3 register.

  @Preconditions
    The TMR3_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    TMR3_Initialize();
    TMR3_StartTimer();

    if(TMR3_HasOverflowOccured())
    {
        TMR3_StopTimer();
    }

    TMR3_Reload();}
    </code>
*/
void TMR3_Reload(void);

/**
  @Summary
    Starts the single pulse acquisition in TMR3 gate operation.

  @Description
    This function starts the single pulse acquisition in TMR3 gate operation.
    This function must be used when the TMR3 gate is enabled.

  @Preconditions
    Initialize  the TMR3 with gate enable before calling this function.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    uint16_t xVal;
    uint16_t yVal;

    //enable TMR3 singlepulse mode
    TMR3_StartSinglePulseAcquistion();

    //check TMR3 gate status
    if(TMR3_CheckGateValueStatus()== 0)
    {
        xVal = TMR3_ReadTimer();
    }

    // wait untill gate interrupt occured
    while(TMR3GIF == 0)
    {
    }

    yVal = TMR3_ReadTimer();
    </code>
*/
void TMR3_StartSinglePulseAcquisition(void);

/**
  @Summary
    Check the current state of Timer1 gate.

  @Description
    This function reads the TMR3 gate value and return it.
    This function must be used when the TMR3 gate is enabled.

  @Preconditions
    Initialize  the TMR3 with gate enable before calling this function.

  @Param
    None

  @Returns
    None

  @Example
    <code>
    uint16_t xVal;
    uint16_t yVal;

    //enable TMR3 singlepulse mode
    TMR3_StartSinglePulseAcquistion();

    //check TMR3 gate status
    if(TMR3_CheckGateValueStatus()== 0)
    {
        xVal = TMR3_ReadTimer();
    }

    //wait untill gate interrupt occured
    while(TMR3IF == 0)
    {
    }

    yVal = TMR3_ReadTimer();
    </code>
*/
uint8_t TMR3_CheckGateValueStatus(void);

/**
  @Summary
    Get the TMR3 overflow status.

  @Description
    This routine get the TMR3 overflow status.

  @Preconditions
    The TMR3_Initialize() routine should be called
    prior to use this routine.

  @Param
    None

  @Returns
    true  - Overflow has occured.
    false - Overflow has not occured.

  @Comment


  @Example
     <code>
     TMR3_Initialize();

     TMR3_StartTimer();

    while(1)
    {
        if(TMR3_HasOverflowOccured())
        {
            TMR3_StopTimer();
        }
    }
    </code>
*/
bool TMR3_HasOverflowOccured(void);


#ifdef __cplusplus  // Provide C++ Compatibility

    }

#endif

#endif // _TMR3_H
/**
 End of File
*/