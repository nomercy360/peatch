import { User } from "./User";

 export type UserWithToken = {
    /**
     * @type array | undefined
    */
    following?: number[];
    /**
     * @type string | undefined
    */
    token?: string;
    user?: User;
};