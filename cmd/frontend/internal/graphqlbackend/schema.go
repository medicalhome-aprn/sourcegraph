package graphqlbackend

var Schema = `schema {
	query: Query
	mutation: Mutation
}

type EmptyResponse {
	alwaysNil: String
}

interface Node {
	id: ID!
}

type Query {
	root: Root!
	node(id: ID!): Node
}

type Mutation {
	createThread(remoteURI: String!, accessToken: String!, file: String!, revision: String!, startLine: Int!, endLine: Int!, startCharacter: Int!, endCharacter: Int!, contents: String!, authorName: String!, authorEmail: String!): Thread!
	createThread2(orgID: Int!, remoteURI: String!, file: String!, revision: String!, startLine: Int!, endLine: Int!, startCharacter: Int!, endCharacter: Int!, contents: String!): Thread!
	updateThread(remoteURI: String!, accessToken: String!, threadID: Int!, archived: Boolean): Thread!
	updateThread2(threadID: Int!, archived: Boolean): Thread!
	addCommentToThread(threadID: Int!, remoteURI: String!, accessToken: String!, contents: String!, authorName: String!, authorEmail: String!): Thread!
	addCommentToThread2(threadID: Int!, contents: String!): Thread!
	createOrg(name: String!, username: String!, email: String!, displayName: String!, avatarUrl: String!): Org!
	inviteUser(email: String!, orgID: Int!): EmptyResponse
	acceptUserInvite(inviteToken: String!, username: String!, email: String!, displayName: String!, avatarUrl: String!): OrgMember!
	removeUserFromOrg(userID: String!, orgID: Int!): EmptyResponse
}


type Root {
	repository(uri: String!): Repository
	repositories(query: String = ""): [Repository!]!
	symbols(id: String!, mode: String!): [Symbol!]!
	currentUser: User
	activeRepos(): ActiveRepoResults!
	search(query: String = "", repositories: [String!]!, first: Int): [SearchResult!]!
	searchRepos(query: SearchQuery!, repositories: [RepositoryRevision!]!): SearchResults!
	searchProfiles: [SearchProfile!]!
	revealCustomerCompany(ip: String!): CompanyProfile
	threads(remoteURI: String!, accessToken: String!, file: String, limit: Int): [Thread!]!
	org(id: Int!): Org!
	packages(lang: String!, id: String, type: String, name: String, commit: String, baseDir: String, repoURL: String, version: String, limit: Int): [Package!]!
	dependents(lang: String!, id: String, type: String, name: String, commit: String, baseDir: String, repoURL: String, version: String, package: String, limit: Int): [Dependency!]!
}

union SearchResult = Repository | File | SearchProfile

type RefFields {
	refLocation: RefLocation
	uri: URI
}

type URI {
	host: String!
	fragment: String!
	path: String!
	query: String!
	scheme: String!
}

type RefLocation {
	startLineNumber: Int!
	startColumn: Int!
	endLineNumber: Int!
	endColumn: Int!
}

type Repository implements Node {
	id: ID!
	uri: String!
	description: String!
	language: String!
	fork: Boolean!
	starsCount: Int
	forksCount: Int
	private: Boolean!
	createdAt: String!
	pushedAt: String!
	commit(rev: String!): CommitState!
	revState(rev: String!): RevState!
	latest: CommitState!
	lastIndexedRevOrLatest: CommitState!
	defaultBranch: String!
	branches: [String!]!
	tags: [String!]!
	listTotalRefs: TotalRefList!
	gitCmdRaw(params: [String!]!): String!
}

type TotalRefList {
	repositories: [Repository!]!
	total: Int!
}

type Symbol {
	repository: Repository!
	path: String!
	line: Int!
	character: Int!
}

type CommitState {
	commit: Commit
	cloneInProgress: Boolean!
}

type RevState {
	commit: Commit
	cloneInProgress: Boolean!
}

input SearchQuery {
	pattern: String!
	isRegExp: Boolean!
	isWordMatch: Boolean!
	isCaseSensitive: Boolean!
	fileMatchLimit: Int!
	includePattern: String
	excludePattern: String
}

input RepositoryRevision {
	repo: String!
	rev: String
}

type Commit implements Node {
	id: ID!
	sha1: String!
	tree(path: String = "", recursive: Boolean = false): Tree
	file(path: String!): File
	languages: [String!]!
}

type CommitInfo {
	rev: String!
	author: Signature
	committer: Signature
	message: String!
}

type Signature {
	person: Person
	date: String!
}

type Person {
	name:  String!
	email: String!
	gravatarHash: String!
}

type Tree {
	directories: [Directory]!
	files: [File]!
}

type Directory {
	name: String!
	tree: Tree!
}

type HighlightedFile {
	aborted: Boolean!
	html: String!
}

type File {
	name: String!
	content: String!
	binary: Boolean!
	isDirectory: Boolean!
	highlight(disableTimeout: Boolean!): HighlightedFile!
	blame(startLine: Int!, endLine: Int!): [Hunk!]!
	commits: [CommitInfo!]!
	dependencyReferences(Language: String!, Line: Int!, Character: Int!): DependencyReferences!
	blameRaw(startLine: Int!, endLine: Int!): String!
}

type ActiveRepoResults {
	active: [String!]!
	inactive: [String!]!
}

type SearchProfile {
	name: String!
	description: String
	repositories: [Repository!]!
}

type SearchResults {
	results: [FileMatch!]!
	limitHit: Boolean!
	cloning: [String!]!
	missing: [String!]!
}

type FileMatch {
	resource: String!
	lineMatches: [LineMatch!]!
	limitHit: Boolean!
}

type LineMatch {
	preview: String!
	lineNumber: Int!
	offsetAndLengths: [[Int!]!]!
	limitHit: Boolean!
}

type DependencyReferences {
	dependencyReferenceData: DependencyReferencesData!
	repoData: RepoDataMap!
}

type RepoDataMap {
	repos: [Repository!]!
	repoIds: [Int!]!
}

type DependencyReferencesData {
	references: [DependencyReference!]!
	location: DepLocation!
}

type DependencyReference {
	dependencyData: String!
	repoId: Int!
	hints: String!
}

type DepLocation {
	location: String!
	symbol: String!
}

type Hunk {
	startLine: Int!
	endLine: Int!
	startByte: Int!
	endByte: Int!
	rev: String!
	author: Signature
	message: String!
}

type Installation {
	login: String!
	githubId: Int!
	installId: Int!
	type: String!
	avatarURL: String!
}

type User {
	githubInstallations: [Installation!]!
	id: String!
	avatarURL: String
	email: String
	orgs: [Org!]!
	orgMemberships: [OrgMember!]!
}

type CompanyProfile {
	ip: String!
	domain: String!
	fuzzy: Boolean!
	company: CompanyInfo!
}

type CompanyInfo {
	id: String!
	name: String!
	legalName: String!
	domain: String!
	domainAliases: [String!]!
	url: String!
	site: SiteDetails!
	category: CompanyCategory!
	tags: [String!]!
	description: String!
	foundedYear: String!
	location: String!
	logo: String!
	tech: [String!]!
}

type SiteDetails {
	url: String!
	title: String!
	phoneNumbers: [String!]!
	emailAddresses: [String!]!
}

type CompanyCategory {
	sector: String!
	industryGroup: String!
	industry: String!
	subIndustry: String!
}

type Org {
	id: Int!
	name: String!
	members: [OrgMember!]!
	repos: [OrgRepo!]!
	threads(file: String, limit: Int): [Thread!]!
}

type OrgMember {
	id: Int!
	org: Org!
	userID: String!
	username: String!
	email: String!
	displayName: String!
	avatarURL: String!
	createdAt: String!
	updatedAt: String!
}

type OrgRepo {
	id: Int!
	org: Org!
	remoteUri: String!
	createdAt: String!
	updatedAt: String!
}

type Thread {
	id: Int!
	repo: OrgRepo!
	file: String!
	revision: String!
	title: String!
	startLine: Int!
	endLine: Int!
	startCharacter: Int!
	endCharacter: Int!
	createdAt: String!
	archivedAt: String
	comments: [Comment!]!
}

type Comment {
	id: Int!
	contents: String!
	createdAt: String!
	updatedAt: String!
	authorName: String!
	authorEmail: String!
	author: OrgMember!
}

type Package {
	lang: String!
	repo: Repository

	# The following fields are properties of build package configuration as returned by the workspace/xpackages LSP endpoint.
	id: String
	type: String
	name: String
	commit: String
	baseDir: String
	repoURL: String
	version: String
}

type Dependency {
	repo: Repository

	# The following fields are properties of build package configuration as returned by the workspace/xpackages LSP endpoint.
	name: String
	repoURL: String
	depth: Int
	vendor: Boolean
	package: String
	absolute: String
	type: String
	commit: String
	version: String
	id: String
	package: String
}
`
