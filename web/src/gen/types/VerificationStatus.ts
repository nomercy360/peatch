/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

export const verificationStatus = {
  VerificationStatusPending: 'pending',
  VerificationStatusVerified: 'verified',
  VerificationStatusDenied: 'denied',
  VerificationStatusBlocked: 'blocked',
  VerificationStatusUnverified: 'unverified',
} as const

export type VerificationStatusEnum =
  (typeof verificationStatus)[keyof typeof verificationStatus]

export type VerificationStatus = VerificationStatusEnum
