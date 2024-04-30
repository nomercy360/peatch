import type { User } from "./User";

 export type GetApiUsersIdPathParams = {
    /**
     * @description User ID
     * @type integer
    */
    id: number;
};

 /**
 * @description OK
*/
export type GetApiUsersId200 = User;

 /**
 * @description OK
*/
export type GetApiUsersIdQueryResponse = User;

 export type GetApiUsersIdQuery = {
    Response: GetApiUsersIdQueryResponse;
    PathParams: GetApiUsersIdPathParams;
};