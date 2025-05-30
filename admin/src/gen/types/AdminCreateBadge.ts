/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { Badge } from './Badge.ts'
import type { CreateBadgeRequest } from './CreateBadgeRequest.ts'

/**
 * @description Created
 */
export type AdminCreateBadge201 = Badge

/**
 * @description Badge data
 */
export type AdminCreateBadgeMutationRequest = CreateBadgeRequest

export type AdminCreateBadgeMutationResponse = AdminCreateBadge201

export type AdminCreateBadgeMutation = {
  Response: AdminCreateBadge201
  Request: AdminCreateBadgeMutationRequest
  Errors: any
}