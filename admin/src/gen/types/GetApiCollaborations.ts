/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { CollaborationResponse } from './CollaborationResponse.ts'

export type GetApiCollaborationsQueryParams = {
  /**
   * @description Page
   * @type integer | undefined
   */
  page?: number
  /**
   * @description Limit
   * @type integer | undefined
   */
  limit?: number
  /**
   * @description Order by
   * @type string | undefined
   */
  order?: string
}

/**
 * @description OK
 */
export type GetApiCollaborations200 = CollaborationResponse[]

export type GetApiCollaborationsQueryResponse = GetApiCollaborations200

export type GetApiCollaborationsQuery = {
  Response: GetApiCollaborations200
  QueryParams: GetApiCollaborationsQueryParams
  Errors: any
}