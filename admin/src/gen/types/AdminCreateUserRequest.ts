/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { Link } from './Link.ts'

export type AdminCreateUserRequest = {
  /**
   * @type array | undefined
   */
  badges?: string[]
  /**
   * @type integer | undefined
   */
  chat_id?: number
  /**
   * @type string | undefined
   */
  description?: string
  /**
   * @type array | undefined
   */
  links?: Link[]
  /**
   * @type string | undefined
   */
  location?: string
  /**
   * @type string | undefined
   */
  name?: string
  /**
   * @type array | undefined
   */
  opportunity_ids?: string[]
  /**
   * @type string | undefined
   */
  title?: string
  /**
   * @type string | undefined
   */
  username?: string
}