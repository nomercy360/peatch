import type { User } from "./User";

 export type GetApiUsersQueryParams = {
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
    /**
     * @description Search
     * @type string | undefined
    */
    search?: string;
    /**
     * @description Find similar
     * @type boolean | undefined
    */
    find_similar?: boolean;
};

 /**
 * @description OK
*/
export type GetApiUsers200 = User[];

 /**
 * @description OK
*/
export type GetApiUsersQueryResponse = User[];

 export type GetApiUsersQuery = {
    Response: GetApiUsersQueryResponse;
    QueryParams: GetApiUsersQueryParams;
};