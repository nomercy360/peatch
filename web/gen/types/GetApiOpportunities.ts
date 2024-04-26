import type { Opportunity } from "./Opportunity";

 /**
 * @description OK
*/
export type GetApiOpportunities200 = Opportunity[];

 /**
 * @description OK
*/
export type GetApiOpportunitiesQueryResponse = Opportunity[];

 export type GetApiOpportunitiesQuery = {
    Response: GetApiOpportunitiesQueryResponse;
};