describe('Add Secret', function () {
    it('Login', function () {
        cy.visit('http://localhost:9002/webapp/login');

        cy.get('#email').type('admin@email.com').should('have.value', 'admin@email.com');
        cy.get('#password').type('Hello123!').should('have.value', 'Hello123!');
        cy.contains('button', 'Login').should('exist', { timeout: 6000 }).click();
        cy.url().should('include', '/webapp');
    });

    it('Add Secret', function () {
        cy.contains('Secret').should('exist', { timeout: 6000 }).click();
        cy.contains('button', 'Add').should('exist', { timeout: 6000 }).click();

        cy.get('#name').type('CySecret').should('have.value', 'CySecret');
        cy.get('#description').type('Description').should('have.value', 'Description');
        cy.get('#secret').type('Hello123!').should('have.value', 'Hello123!');

        cy.intercept('POST', '/app/private/graphql').as('post');
        cy.contains('Save').should('exist', { timeout: 6000 }).click();

        cy.wait('@post').its('response.body.errors').should('not.exist');
        cy.wait('@post').its('response.statusCode').should('eq', 200);
    });

    it('Delete Secret', function () {
        cy.contains('CySecret').click();

        cy.contains('Delete secret').should('exist', { timeout: 6000 }).click();
        cy.intercept('POST', '/app/private/graphql').as('post');
        cy.contains('Yes').should('exist', { timeout: 6000 }).click();

        cy.wait('@post').its('response.body.errors').should('not.exist');
        cy.wait('@post').its('response.statusCode').should('eq', 200);
    });
});
