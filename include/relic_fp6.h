/*
 * RELIC is an Efficient LIbrary for Cryptography
 * Copyright (C) 2007, 2008, 2009 RELIC Authors
 *
 * This file is part of RELIC. RELIC is legal property of its developers,
 * whose names are not listed here. Please refer to the COPYRIGHT file
 * for contact information.
 *
 * RELIC is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * RELIC is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with RELIC. If not, see <http://www.gnu.org/licenses/>.
 */

/**
 * @defgroup fp6 Sextic extension prime field arithmetic.
 */

/**
 * @file
 *
 * Interface of the sextic extension prime field arithmetic functions.
 *
 * @version $Id$
 * @ingroup fp6
 */

#ifndef RELIC_FP6_H
#define RELIC_FP6_H

#include "relic_fp2.h"
#include "relic_types.h"

/**
 * Represents a sextic extension prime field element.
 * 
 * An element a is represented as a[0] + a[1] * Y + a[1] * Y^2,
 * where Y^3 = X if p = 5 mod 8, Y^3 = 1 + X if p = 3 mod 8 and
 * Y^3 = 2 + X if p = 7 mod 8 and p = 2,3 mod 5.
 */
typedef fp2_t fp6_t[3];

/**
 * Allocate and initializes a sextic extension prime field element.
 *
 * @param[out] A			- the new sextic extension field element.
 */
#define fp6_new(A)															\
		fp2_new(A[0]); fp2_new(A[1]); fp2_new(A[2]);						\

/**
 * Calls a function to clean and free a sextic extension field element.
 *
 * @param[out] A			- the sextic extension field element to free.
 */
#define fp6_free(A)															\
		fp2_free(A[0]); fp2_free(A[1]); fp2_new(A[2]); 						\

/**
 * Copies the second argument to the first argument.
 *
 * @param[out] C			- the result.
 * @param[in] A				- the sextic extension field element to copy.
 */
#define fp6_copy(C, A)														\
		fp2_copy(C[0], A[0]); fp2_copy(C[1], A[1]); fp2_copy(C[2], A[2]);	\

/**
 * Negates a sextic extension field element.
 *
 * @param[out] c			- the result.
 * @param[out] a			- the sextic extension field element to negate.
 */
#define fp6_neg(C, A)														\
		fp2_neg(C[0], A[0]); fp2_neg(C[1], A[1]); fp2_neg(C[2], A[2]); 		\

/**
 * Assigns zero to a sextic extension field element.
 *
 * @param[out] a			- the sextic extension field element to zero.
 */
#define fp6_zero(A)															\
		fp2_zero(A[0]); fp2_zero(A[1]); fp2_zero(A[2]); 					\

/**
 * Tests if a sextic extension field element is zero or not.
 *
 * @param[in] a				- the sextic extension field element to compare.
 * @return 1 if the argument is zero, 0 otherwise.
 */
#define fp6_is_zero(A)														\
		fp2_is_zero(A[0]) || fp2_is_zero(A[1]) || fp2_is_zero(A[2])			\

/**
 * Assigns a random value to a sextic extension field element.
 *
 * @param[out] a			- the sextic extension field element to assign.
 */
#define fp6_rand(A)															\
		fp2_rand(A[0]); fp2_rand(A[1]); fp2_rand(A[2]);						\

/**
 * Prints a sextic extension field element to standard output.
 *
 * @param[in] a				- the sextic extension field element to print.
 */
#define fp6_print(A)														\
		fp2_print(A[0]); fp2_print(A[1]); fp2_print(A[2]);					\

/**
 * Reads a sextic extension field element from strings in a given radix.
 * The radix must be a power of 2 included in the interval [2, 64].
 *
 * @param[out] a			- the result.
 * @param[in] str00			- 
 * @param[in] str01			- 
 * @param[in] str10			- 
 * @param[in] str11			- 
 * @param[in] str20			- 
 * @param[in] str21			- 
 * @param[in] len 			- the maximum length of the strings.
 * @param[in] radix			- the radix.
 * @throw ERR_INVALID		- if the radix is invalid.
 */
#define fp6_read(A, STR00, STR01, STR10, STR11, STR20, STR21, LEN, RADIX)	\
		fp2_read(A[0], STR00, STR01, LEN, RADIX);							\
		fp2_read(A[1], STR10, STR11, LEN, RADIX);							\
		fp2_read(A[2], STR20, STR21, LEN, RADIX);							\

/**
 * Writes a sextic extension field element to strings in a given radix.
 * The radix must be a power of 2 included in the interval [2, 64].
 *
 * @param[out] str00		- 
 * @param[out] str01		- 
 * @param[out] str10		- 
 * @param[out] str11		- 
 * @param[out] str20		- 
 * @param[out] str21		- 
 * @param[in] len			- the buffer capacities.
 * @param[in] a				- the sextic extension field element to write.
 * @param[in] radix			- the radix.
 * @throw ERR_BUFFER		- if the buffer capacity is insufficient.
 * @throw ERR_INVALID		- if the radix is invalid.
 */
#define fp6_write(STR00, STR01, STR10, STR11, STR20, STR21, LEN, A, RADIX)	\
		fp2_write(STR00, STR01, LEN, A[0], RADIX);							\
		fp2_write(STR10, STR11, LEN, A[1], RADIX);							\
		fp2_write(STR20, STR21, LEN, A[2], RADIX);							\

/**
 * Returns the result of a comparison between two sextic extension field
 * elements
 *
 * @param[in] a				- the first sextic extension field element.
 * @param[in] b				- the second sextic extension field element.
 * @return FP_LT if a < b, FP_EQ if a == b and FP_GT if a > b.
 */
#define fp6_cmp(A, B)														\
		((fp2_cmp(A[0], B[0]) == CMP_EQ) && (fp2_cmp(A[1], B[1]) == CMP_EQ)	\
		&& (fp2_cmp(A[2], B[2]) == CMP_EQ) ? CMP_EQ : CMP_NE)				\

#define fp6_set_dig(A, B)													\
		fp2_set_dig(A[0], B); fp2_zero(A[1]); fp2_zero(A[2]);				\

/**
 * Adds two sextic extension field elements, that is, computes c = a + b.
 *
 * @param[out] c			- the destination.
 * @param[in] a				- the first sextic extension field element.
 * @param[in] b				- the second sextic extension field element.
 */
#define fp6_add(C, A, B)													\
		fp2_add(C[0], A[0], B[0]); fp2_add(C[1], A[1], B[1]);				\
		fp2_add(C[2], A[2], B[2]);											\

#define fp6_dbl(C, A)														\
		fp2_dbl(C[0], A[0]); fp2_dbl(C[1], A[1]); fp2_dbl(C[2], A[2]);		\

#define fp6_add_dig(C, A, B)												\
		fp2_add_dig(C[0], A[0], B); fp2_copy(C[1], A[1]);					\
		fp2_copy(C[2], A[2]);												\

#define fp6_sub_dig(C, A, B)												\
		fp2_sub_dig(C[0], A[0], B); fp2_copy(C[1], A[1]);					\
		fp2_copy(C[2], A[2]);

/**
 * Subtracts a sextic extension field element from another, that is, computes
 * c = a - b.
 *
 * @param[out] c			- the destination.
 * @param[in] a				- the sextic extension prime field element.
 * @param[in] b				- the sextic extension prime field element.
 */
#define fp6_sub(C, A, B)													\
		fp2_sub(C[0], A[0], B[0]); fp2_sub(C[1], A[1], B[1]);				\
		fp2_sub(C[2], A[2], B[2]);											\

/**
 * Multiples two sextic extension field elements, that is, compute c = a * b.
 *
 * @param[out] c			- the destination.
 * @param[in] a				- the sextic extension prime field element.
 * @param[in] b				- the sextic extension prime field element.
 */
void fp6_mul(fp6_t c, fp6_t a, fp6_t b);

/**
 * Computes the square of a sextic extension field element, that is, computes
 * c = a * a.
 *
 * @param[out] c			- the destination.
 * @param[in] a				- the sextic extension field element to square.
 */
void fp6_sqr(fp6_t c, fp6_t a);

/**
 * Inverts a sextic extension field element. Computes c = a^(-1).
 *
 * @param[out] c			- the destination.
 * @param[in] a				- the sextic extension prime field element to invert.
 */
void fp6_inv(fp6_t c, fp6_t a);

/**
 * Multiples a sextic extension field element by a quadratic extension field element.
 *
 * @param[out] c			- the destination.
 * @param[in] a				- the sextic extension prime field element.
 * @param[in] b				- the quadratic extension prime field element.
 */
void fp6_mul_fp2(fp6_t c, fp6_t a, fp2_t b);

/**
 * Multiples a sextic extension field element by a sparse element.
 * 
 * The sparse element must have a[2] = 0.
 *
 * @param[out] c			- the destination.
 * @param[in] a				- a sextic extension prime field element.
 * @param[in] b				- a sparse sextic extension prime field element.
 */
void fp6_mul_sparse(fp6_t c, fp6_t a, fp6_t b);

/**
 * Multiplies a sextic extension field element by Y.
 * 
 * @param[out] c			- the destination.
 * @param[in] a				- the sextic extension prime field element.
 */
void fp6_mul_poly(fp6_t c, fp6_t a);

/**
 * Computes the Frobenius endomorphism of a unitary sextic extension field element,
 * that is, Frob(a) = a^p.
 * 
 * @param[out] c			- the destination.
 * @param[in] a				- a sextic extension prime field element.
 * @param[in] b				- constant used in Frobenius, (Z^p)^2 = Y^p. 
 */
void fp6_frob(fp6_t c, fp6_t a, fp6_t b);

#endif /* !RELIC_FP6_H */
