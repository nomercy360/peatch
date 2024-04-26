import type { Collaboration } from "./Collaboration";

 export type GetApiCollaborationsQueryParams = {
    /**
     * @description Page
     * @type integer | undefined
    */
    page?: number;
    /**
     * @description Limit
     * @type integer | undefined
    */
    limit?: number;
    /**
     * @description Order by
     * @type string | undefined
    */
    order?: string;
};

 /**
 * @description OK
*/
export type GetApiCollaborations200 = Collaboration[];

 /**
 * @description OK
*/
export type GetApiCollaborationsQueryResponse = Collaboration[];

 export type GetApiCollaborationsQuery = {
    Response: GetApiCollaborationsQueryResponse;
    QueryParams: GetApiCollaborationsQueryParams;
};