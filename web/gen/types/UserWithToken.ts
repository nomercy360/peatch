import { User } from './User';

export type UserWithToken = {
    /**
     * @description TODO: implement following
     * @type array | undefined
    */
    following?: number[];
    /**
     * @type string | undefined
    */
    token?: string;
    user?: User;
};