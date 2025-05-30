/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { BadgeResponse } from './BadgeResponse.ts'
import type { CityResponse } from './CityResponse.ts'
import type { Link } from './Link.ts'
import type { OpportunityResponse } from './OpportunityResponse.ts'

export type UserProfileResponse = {
  /**
   * @type string | undefined
   */
  avatar_url?: string
  /**
   * @type array | undefined
   */
  badges?: BadgeResponse[]
  /**
   * @type string | undefined
   */
  description?: string
  /**
   * @type string | undefined
   */
  id?: string
  /**
   * @type boolean | undefined
   */
  is_following?: boolean
  /**
   * @type string | undefined
   */
  last_active_at?: string
  /**
   * @type array | undefined
   */
  links?: Link[]
  /**
   * @type object | undefined
   */
  location?: CityResponse
  /**
   * @type string | undefined
   */
  name?: string
  /**
   * @type array | undefined
   */
  opportunities?: OpportunityResponse[]
  /**
   * @type string | undefined
   */
  title?: string
  /**
   * @type string | undefined
   */
  username?: string
}