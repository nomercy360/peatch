import type { Badge } from "./Badge";

 export type GetApiBadgesIdPathParams = {
    /**
     * @description Badge ID
     * @type integer
    */
    id: number;
};

 /**
 * @description OK
*/
export type GetApiBadgesId200 = Badge;

 /**
 * @description OK
*/
export type GetApiBadgesIdQueryResponse = Badge;

 export type GetApiBadgesIdQuery = {
    Response: GetApiBadgesIdQueryResponse;
    PathParams: GetApiBadgesIdPathParams;
};