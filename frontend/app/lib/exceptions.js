export class AuthRequiredError extends Error {
    constructor(message = 'Auth required to access this page') {
        super(message);
        this.name = 'AuthRequiredError';
    }
}