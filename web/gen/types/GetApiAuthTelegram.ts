import type { User } from "./User";

 export type GetApiAuthTelegramQueryParams = {
    /**
     * @description Query ID
     * @type string
    */
    query_id: string;
    /**
     * @description User
     * @type string
    */
    user: string;
    /**
     * @description Auth date
     * @type string
    */
    auth_date: string;
    /**
     * @description Hash
     * @type string
    */
    hash: string;
};

 /**
 * @description OK
*/
export type GetApiAuthTelegram200 = User;

 /**
 * @description OK
*/
export type GetApiAuthTelegramQueryResponse = User;

 export type GetApiAuthTelegramQuery = {
    Response: GetApiAuthTelegramQueryResponse;
    QueryParams: GetApiAuthTelegramQueryParams;
};