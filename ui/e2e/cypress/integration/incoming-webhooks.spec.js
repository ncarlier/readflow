/// <reference types="Cypress" />

const webhookAlias = 'automated-test'

context('Settings - Incoming Webhooks', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/login')
  })

  it('should create an incoming webhook', () => {
    // Go to integration
    cy.get('#navbar-link-settings').click()
    cy.get('a[href="/settings/integrations"]').click()
    cy.title().should('include', 'Integrations')
    // Add an new incoming webhook
    cy.get('#add-new-incoming-webhook').click()
    cy.title().should('include', 'Add new incoming webhook')
    cy.focused().should('have.attr', 'name', 'alias').type(webhookAlias)
    cy.get('[data-test="btn-primary"]').contains('Add').click()
    // Validate incoming webhook creation
    cy.title().should('include', 'Integrations')
    cy.get(`[data-test-id="${webhookAlias}"]`).contains(webhookAlias)
  })

  it('should udate an incoming webhook', () => {
    cy.get('#navbar-link-settings').click()
    cy.get('a[href="/settings/integrations"]').click()
    cy.title().should('include', 'Integrations')
    cy.get(`[data-test-id="${webhookAlias}"]`).contains(webhookAlias).click()
    cy.title().should('include', 'Edit incoming Webhook')
    cy.focused().should('have.attr', 'name', 'alias').type('-edited')
    cy.get('[data-test="btn-primary"]').contains('Update').click()
    cy.title().should('include', 'Integrations')
    cy.get(`[data-test-id="${webhookAlias}-edited"]`).contains(`${webhookAlias}-edited`)
  })

  it('should delete an incoming webhook', () => {
    cy.get('#navbar-link-settings').click()
    cy.get('a[href="/settings/integrations"]').click()
    cy.title().should('include', 'Integrations')
    cy.get(`[data-test-id="${webhookAlias}-edited"]`).closest('tr').find('input[type=checkbox]').click()
    cy.get('#remove-selection-1').click()
    cy.get('.ReactModal__Content').should('be.visible')
    cy.get('.ReactModal__Content header > h1').contains('Delete incoming webhook?')
    cy.get('.ReactModal__Content [data-test="btn-primary"]').contains('Delete').click()
    cy.get(`[data-test-id="${webhookAlias}-edited"]`).should('not.exist')
  })
})
