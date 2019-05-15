/// <reference types="Cypress" />

context('Viewport', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/login')
  })

  it('should create an API key', () => {
    // Go to API keys page
    cy.get('#navbar-link-settings').click()
    cy.get('a[href="/settings/api-keys"]').click()
    cy.title().should('include', 'API keys')
    // Add an new API key
    cy.get('#add-new-api-key').click()
    cy.title().should('include', 'Add new API key')
    cy.focused()
      .should('have.attr', 'name', 'alias')
      .type('test')
    cy.get('[data-test="btn-primary"]')
      .contains('Add')
      .click()
    // Validate API key creation
    cy.title().should('include', 'API keys')
    cy.get('[data-test-id="test"] a').contains('test')
  })

  it('should udate an API key', () => {
    cy.get('#navbar-link-settings').click()
    cy.get('a[href="/settings/api-keys"]').click()
    cy.title().should('include', 'API keys')
    cy.get('[data-test-id="test"] a')
      .contains('test')
      .click()
    cy.title().should('include', 'Edit API key')
    cy.focused()
      .should('have.attr', 'name', 'alias')
      .type('-test')
    cy.get('[data-test="btn-primary"]')
      .contains('Update')
      .click()
    cy.title().should('include', 'API keys')
    cy.get('[data-test-id="test-test"] a').contains('test-test')
  })

  it('should delete an API key', () => {
    cy.get('#navbar-link-settings').click()
    cy.get('a[href="/settings/api-keys"]').click()
    cy.title().should('include', 'API keys')
    cy.get('[data-test-id="test-test"] input[type=checkbox]').click()
    cy.get('#remove-selection').click()
    cy.get('.ReactModal__Content').should('be.visible')
    cy.get('.ReactModal__Content header > h1').contains('Delete API keys?')
    cy.get('.ReactModal__Content [data-test="btn-primary"]')
      .contains('Delete')
      .click()
    cy.get('[data-test-id="test-test]"').should('not.exist')
  })
})
