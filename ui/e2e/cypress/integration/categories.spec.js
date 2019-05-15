/// <reference types="Cypress" />

context('Viewport', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/login')
  })

  it('should create a category', () => {
    // Go to categories page
    cy.get('#navbar-link-settings').click()
    cy.title().should('include', 'Categories')
    // Add an new category
    cy.get('#add-new-category').click()
    cy.title().should('include', 'Add new category')
    cy.focused()
      .should('have.attr', 'name', 'title')
      .type('test')
    cy.get('[data-test="btn-primary"]')
      .contains('Add')
      .click()
    // Validate category creation
    cy.title().should('include', 'Categories')
    cy.get('[data-test-id="test"] a').contains('test')
  })

  it('should udate a category', () => {
    // Go to categories page
    cy.get('#navbar-link-settings').click()
    cy.title().should('include', 'Categories')
    cy.get('[data-test-id="test"] a')
      .contains('test')
      .click()
    cy.title().should('include', 'Edit category')
    cy.focused()
      .should('have.attr', 'name', 'title')
      .type('-test')
    cy.get('[data-test="btn-primary"]')
      .contains('Update')
      .click()
    cy.title().should('include', 'Categories')
    cy.get('[data-test-id="test-test"] a').contains('test-test')
  })

  it('should delete a category', () => {
    // Go to categories page
    cy.get('#navbar-link-settings').click()
    cy.title().should('include', 'Categories')
    cy.get('[data-test-id="test-test"] input[type=checkbox]').click()
    cy.get('#remove-selection').click()
    cy.get('.ReactModal__Content').should('be.visible')
    cy.get('.ReactModal__Content header > h1').contains('Delete categories?')
    cy.get('.ReactModal__Content [data-test="btn-primary"]')
      .contains('Delete')
      .click()
    cy.get('[data-test-id="test-test]"').should('not.exist')
  })
})
