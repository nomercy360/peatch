/**
 * Generated by Kubb (https://kubb.dev/).
 * Do not edit manually.
 */

import type { Opportunity } from './Opportunity.ts'

/**
 * @description OK
 */
export type AdminListOpportunities200 = Opportunity[]

export type AdminListOpportunitiesQueryResponse = AdminListOpportunities200

export type AdminListOpportunitiesQuery = {
  Response: AdminListOpportunities200
  Errors: any
}