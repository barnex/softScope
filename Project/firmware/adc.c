#include <stm32f4xx.h>
#include <stdint.h>
#include <string.h>

#include "adc.h"

void init_analogIn() {
	GPIO_InitTypeDef gpio = {0, };
	RCC_AHB1PeriphClockCmd(RCC_AHB1Periph_GPIOA, ENABLE);
	memset((void*) &gpio, 0, sizeof(GPIO_InitTypeDef));
	gpio.GPIO_Pin = GPIO_Pin_1;
	gpio.GPIO_Mode = GPIO_Mode_AN;
	GPIO_Init(GPIOA, &gpio);
}

#define ADC_DR          0x4C                         // ADC data register offset, STM32 Guide section 13.13.14
#define ADC1_DR_ADDRESS 0x40012000 + 0x000 + ADC_DR  // ADC1 data register, STM32 Guide section 2.3
#define ADC2_DR_ADDRESS 0x40012000 + 0x100 + ADC_DR
#define ADC3_DR_ADDRESS 0x40012000 + 0x200 + ADC_DR


void init_ADC(uint16_t volatile *adc1buf, uint16_t volatile *adc2buf, int samples) {

	// Init the DMA for transferring data from the ADC
	// Enable the clock to the DMA
	RCC_AHB1PeriphClockCmd(RCC_AHB1Periph_DMA2, ENABLE);

	// Configure the DMA stream for ADC -> memory
	DMA_InitTypeDef DMAInit = {0, };
	DMAInit.DMA_Channel            = DMA_Channel_0;                   // DMA channel 0 stream 0 is mapped to ADC1
	DMAInit.DMA_PeripheralBaseAddr = ADC1_DR_ADDRESS;  
	DMAInit.DMA_Memory0BaseAddr	   = (uint32_t)(adc1buf);  
	DMAInit.DMA_DIR	               = DMA_DIR_PeripheralToMemory;
	DMAInit.DMA_BufferSize         = samples;                         // 
	DMAInit.DMA_PeripheralInc      = DMA_PeripheralInc_Disable;       // Do not increase the periph pointer
	DMAInit.DMA_MemoryInc          = DMA_MemoryInc_Enable;            // But do increase the memory pointer
	DMAInit.DMA_PeripheralDataSize = DMA_PeripheralDataSize_HalfWord; //16 bits only please
	DMAInit.DMA_MemoryDataSize	   = DMA_MemoryDataSize_HalfWord;
	DMAInit.DMA_Mode               = DMA_Mode_Circular;               // Wrap around and keep playing a shanty tune
	DMAInit.DMA_Priority           = DMA_Priority_VeryHigh;
	DMAInit.DMA_FIFOMode           = DMA_FIFOMode_Disable;            // No FIFO, direct write will be sufficiently fast
	DMAInit.DMA_MemoryBurst        = DMA_MemoryBurst_Single;
	DMAInit.DMA_PeripheralBurst    = DMA_PeripheralBurst_Single;

	// Initialize the DMA
	DMA_Init(DMA2_Stream0, &DMAInit);
	DMA_Cmd(DMA2_Stream0 , ENABLE);

	DMAInit.DMA_Channel            = DMA_Channel_0;                   // DMA channel 0 stream 1 is mapped to ADC2
	DMAInit.DMA_PeripheralBaseAddr = ADC2_DR_ADDRESS;  
	DMAInit.DMA_Memory0BaseAddr	   = (uint32_t)(adc2buf);   
	DMA_Init(DMA2_Stream1, &DMAInit);
	DMA_Cmd(DMA2_Stream1 , ENABLE);

	//Enable the clock to the ADC
	RCC_APB2PeriphClockCmd(RCC_APB2Periph_ADC1, ENABLE);
	RCC_APB2PeriphClockCmd(RCC_APB2Periph_ADC2, ENABLE);

	// The things that are shared between the three ADCs
	ADC_CommonInitTypeDef common = {0, };
	common.ADC_Mode      = ADC_Mode_Independent;
	common.ADC_Prescaler = ADC_Prescaler_Div2;
	ADC_CommonInit(&common);

	// The things specific to ADC1
	ADC_InitTypeDef adc = {0, };
	adc.ADC_Resolution           = ADC_Resolution_12b;
	adc.ADC_ScanConvMode         = DISABLE; // Disable scanning multiple channels
	adc.ADC_ContinuousConvMode   = DISABLE; // Disable ADC free running
	adc.ADC_ExternalTrigConvEdge = ADC_ExternalTrigConvEdge_Rising;
	adc.ADC_ExternalTrigConv     = ADC_ExternalTrigConv_T2_TRGO; // Trigger of TIM2
	adc.ADC_DataAlign            = ADC_DataAlign_Right;
	adc.ADC_NbrOfConversion      = 1;
	ADC_Init(ADC1, &adc);
	ADC_DMACmd(ADC1, ENABLE); // Enable generating DMA requests
	ADC_DMARequestAfterLastTransferCmd(ADC1, ENABLE);
	ADC_RegularChannelConfig(ADC1, ADC_Channel_1, 1, ADC_SampleTime_3Cycles); // Configure the channel from which to sample

	// The things specific to ADC2
	ADC_Init(ADC2, &adc);
	ADC_DMACmd(ADC2, ENABLE); // Enable generating DMA requests
	ADC_DMARequestAfterLastTransferCmd(ADC2, ENABLE);
	ADC_RegularChannelConfig(ADC2, ADC_Channel_1, 1, ADC_SampleTime_3Cycles); // Configure the channel from which to sample

	// NVIC config
	NVIC_InitTypeDef NVIC_InitStructure;
	NVIC_InitStructure.NVIC_IRQChannel = ADC_IRQn;
	NVIC_InitStructure.NVIC_IRQChannelPreemptionPriority = 1;
	NVIC_InitStructure.NVIC_IRQChannelSubPriority = 0;
	NVIC_InitStructure.NVIC_IRQChannelCmd = ENABLE;
	NVIC_Init(&NVIC_InitStructure);
	ADC_ITConfig(ADC1, ADC_IT_OVR, ENABLE);
	ADC_ClearFlag(ADC1, ADC_FLAG_OVR);
	ADC_ClearITPendingBit(ADC1 , ADC_IT_OVR);

	ADC_Cmd( ADC1, ENABLE );
	ADC_Cmd( ADC2, ENABLE );

	ADC_SoftwareStartConv(ADC1);
	ADC_SoftwareStartConv(ADC2);
}

void ADC_IRQHandler(void) {
	ADC_ClearITPendingBit(ADC1 , ADC_IT_OVR);
	ADC_ClearFlag(ADC1, ADC_FLAG_OVR);
	DMA_Cmd( DMA2_Stream0, ENABLE );
}

