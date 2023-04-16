// +build relic

#ifndef _REL_THRESHOLD_INCLUDE_H
#define _REL_THRESHOLD_INCLUDE_H

#include "bls_include.h"

int G1_lagrangeInterpolateAtZero_serialized(byte*, const byte* , const uint8_t[], const int);
extern void Fr_polynomial_image(Fr* out, G2* y, const Fr* a, const int a_size, const byte x);

#endif
