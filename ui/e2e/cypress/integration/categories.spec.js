/// <reference types="Cypress" />

const categoryName = 'automated-test'

context('Settings - Categories', () => {
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
    cy.focused().should('have.attr', 'name', 'title').type(categoryName)
    cy.get('[data-test="btn-primary"]').contains('Add').click()
    // Validate category creation
    cy.title().should('include', 'Categories')
    cy.get(`[data-test-id="${categoryName}"]`).contains(categoryName)
  })

  it('should udate a category', () => {
    // Go to categories page
    cy.get('#navbar-link-settings').click()
    cy.title().should('include', 'Categories')
    cy.get(`[data-test-id="${categoryName}"]`).contains(categoryName).click()
    cy.title().should('include', 'Edit category')
    cy.focused().should('have.attr', 'name', 'title').type('-edited')
    cy.get('[data-test="btn-primary"]').contains('Update').click()
    cy.title().should('include', 'Categories')
    cy.get(`[data-test-id="${categoryName}-edited"]`).contains(`${categoryName}-edited`)
  })

  it('should delete a category', () => {
    // Go to categories page
    cy.get('#navbar-link-settings').click()
    cy.title().should('include', 'Categories')
    cy.get(`[data-test-id="${categoryName}-edited"]`).closest('tr').find('input[type=checkbox]').click()
    cy.get('#remove-selection').click()
    cy.get('.ReactModal__Content').should('be.visible')
    cy.get('.ReactModal__Content header > h1').contains('Delete categories?')
    cy.get('.ReactModal__Content [data-test="btn-primary"]').contains('Delete').click()
    cy.get(`[data-test-id="${categoryName}-edited]"`).should('not.exist')
  })
})
