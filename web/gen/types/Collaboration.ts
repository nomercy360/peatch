import { Badge } from './Badge';
import { Opportunity } from './Opportunity';

export type Collaboration = {
    /**
     * @type array | undefined
     */
    badges?: Badge[];
  /**
     * @type string | undefined
    */
    city?: string;
    /**
     * @type string | undefined
    */
    country?: string;
    /**
     * @type string | undefined
    */
    country_code?: string;
    /**
     * @type string | undefined
    */
    created_at?: string;
    /**
     * @type string | undefined
    */
    description?: string;
    /**
     * @type integer | undefined
    */
    id?: number;
    /**
     * @type boolean | undefined
    */
    is_payable?: boolean;
  opportunity?: Opportunity;
    /**
     * @type integer | undefined
    */
    opportunity_id?: number;
    /**
     * @type string | undefined
    */
    published_at?: string;
    /**
     * @type integer | undefined
    */
    requests_count?: number;
    /**
     * @type string | undefined
    */
    title?: string;
    /**
     * @type string | undefined
    */
    updated_at?: string;
    /**
     * @type integer | undefined
    */
    user_id?: number;
};