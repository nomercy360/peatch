/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { AdminAuthResponse } from './AdminAuthResponse.ts'
import type { AdminLoginRequest } from './AdminLoginRequest.ts'
import type { ErrorResponse } from './ErrorResponse.ts'

/**
 * @description OK
 */
export type AdminLogin200 = AdminAuthResponse

/**
 * @description Bad Request
 */
export type AdminLogin400 = ErrorResponse

/**
 * @description Admin login credentials
 */
export type AdminLoginMutationRequest = AdminLoginRequest

export type AdminLoginMutationResponse = AdminLogin200

export type AdminLoginMutation = {
  Response: AdminLogin200
  Request: AdminLoginMutationRequest
  Errors: AdminLogin400
}
