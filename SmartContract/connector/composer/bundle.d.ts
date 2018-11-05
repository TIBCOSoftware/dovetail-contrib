

export class AssetDeclaration {
    constructor(arg0?: any, arg1?: any);
    isRelationshipTarget(): any;
    getSystemType(): any;
    validate(): any;
  }
  export class BaseException {
    constructor(arg0?: any);
  }
  export class BaseFileException {
    constructor(arg0?: any, arg1?: any, arg2?: any);
    getFileLocation(): any;
    getShortMessage(): any;
  }
  export class BusinessNetworkDefinition {
    static fromArchive(arg0?: any): any;
    constructor(arg0?: any, arg1?: any, arg2?: any, arg3?: any);
    accept(arg0?: any, arg1?: any): any;
    getModelManager(): any;
  }
  
  export class JSONSchemaVisitor {
    visit(arg0?: any, arg1?: any): any;
    visitBusinessNetwork(arg0?: any, arg1?: any): any;
    visitModelManager(arg0?: any, arg1?: any): any;
    visitModelFile(arg0?: any, arg1?: any): any;
    visitAssetDeclaration(arg0?: any, arg1?: any): any;
    visitTransactionDeclaration(arg0?: any, arg1?: any): any;
    visitConceptDeclaration(arg0?: any, arg1?: any): any;
    visitClassDeclaration(arg0?: any, arg1?: any): any;
    visitClassDeclarationCommon(arg0?: any, arg1?: any, arg2?: any): any;
    visitField(arg0?: any, arg1?: any): any;
    visitEnumDeclaration(arg0?: any, arg1?: any): any;
    visitEnumValueDeclaration(arg0?: any, arg1?: any): any;
    visitRelationshipDeclaration(arg0?: any, arg1?: any): any;
  }
  
  export class ClassDeclaration {
    constructor(arg0?: any, arg1?: any);
    getModelFile(): any;
    process(): any;
    _resolveSuperType(): any;
    validate(): any;
    getSystemType(): any;
    isAbstract(): any;
    isEnum(): any;
    isConcept(): any;
    isEvent(): any;
    isRelationshipTarget(): any;
    isSystemRelationshipTarget(): any;
    isSystemType(): any;
    isSystemCoreType(): any;
    getName(): any;
    getNamespace(): any;
    getFullyQualifiedName(): any;
    getIdentifierFieldName(): any;
    getOwnProperty(arg0?: any): any;
    getOwnProperties(): any;
    getSuperType(): any;
    getSuperTypeDeclaration(): any;
    getAssignableClassDeclarations(): any;
    getAllSuperTypeDeclarations(): any;
    getProperty(arg0?: any): any;
    getProperties(): any;
    getNestedProperty(arg0?: any): any;
    toString(): any;
  }
  export class Concept {
    constructor(arg0?: any, arg1?: any, arg2?: any, arg3?: any);
    isConcept(): any;
  }
  export class ConceptDeclaration {
    constructor(arg0?: any, arg1?: any);
    isConcept(): any;
  }
  
  export class EnumDeclaration {
    constructor(arg0?: any, arg1?: any);
    isEnum(): any;
    toString(): any;
  }
  export class EnumValueDeclaration {
    constructor(arg0?: any, arg1?: any);
    validate(arg0?: any): any;
  }
  export class EventDeclaration {
    constructor(arg0?: any, arg1?: any);
    getSystemType(): any;
    validate(): any;
    isEvent(): any;
  }
  
  export class Field {
    constructor(arg0?: any, arg1?: any);
    process(): any;
    getValidator(): any;
    getDefaultValue(): any;
    toString(): any;
  }
  
  export class Globalize {
    static messageFormatter(arg0?: any): any;
    static formatMessage(arg0?: any): any;
    constructor(arg0?: any);
  }
  
  export class Introspector {
    constructor(arg0?: any);
    accept(arg0?: any, arg1?: any): any;
    getClassDeclarations(): any;
    getClassDeclaration(arg0?: any): any;
    getModelManager(): any;
  }
  
  export class Logger {
    static setFunctionalLogger(arg0?: any): any;
    static getSelectionTree(): any;
    static setLoggerCfg(arg0?: any): any;
    static getLoggerCfg(): any;
    static processLoggerConfig(arg0?: any): any;
    static getLog(arg0?: any): any;
    static _setupLog(arg0?: any): any;
    static _parseLoggerConfig(arg0?: any): any;
    static _loadLogger(arg0?: any): any;
    static setCallBack(arg0?: any): any;
    static getCallBack(): any;
    static __reset(): any;
    static setCLIDefaults(): any;
    static invokeAllLevels(arg0?: any): any;
    constructor(arg0?: any);
    padRight(arg0?: any, arg1?: any): any;
    intlog(arg0?: any, arg1?: any, arg2?: any): any;
    debug(arg0?: any, arg1?: any, arg2?: any): any;
    warn(arg0?: any, arg1?: any, arg2?: any): any;
    info(arg0?: any, arg1?: any, arg2?: any): any;
    verbose(arg0?: any, arg1?: any, arg2?: any): any;
    error(arg0?: any, arg1?: any, arg2?: any): any;
    entry(arg0?: any, arg1?: any): any;
    exit(arg0?: any, arg1?: any): any;
  }
  
  export class ModelFile {
    constructor(arg0?: any, arg1?: any, arg2?: any);
    isExternal(): any;
    getImportURI(arg0?: any): any;
    getExternalImports(): any;
    accept(arg0?: any, arg1?: any): any;
    getModelManager(): any;
    getImports(): any;
    validate(): any;
    resolveType(arg0?: any, arg1?: any): any;
    isLocalType(arg0?: any): any;
    isImportedType(arg0?: any): any;
    resolveImport(arg0?: any): any;
    isDefined(arg0?: any): any;
    getType(arg0?: any): any;
    getFullyQualifiedTypeName(arg0?: any): any;
    getLocalType(arg0?: any): any;
    getAssetDeclaration(arg0?: any): any;
    getTransactionDeclaration(arg0?: any): any;
    getEventDeclaration(arg0?: any): any;
    getParticipantDeclaration(arg0?: any): any;
    getNamespace(): any;
    getName(): any;
    getAssetDeclarations(): any;
    getTransactionDeclarations(): any;
    getEventDeclarations(): any;
    getParticipantDeclarations(): any;
    getConceptDeclarations(): any;
    getEnumDeclarations(): any;
    getDeclarations(arg0?: any): any;
    getAllDeclarations(): any;
    getDefinitions(): any;
    isSystemModelFile(): any;
  }
  export class ModelManager {
    constructor();
    addSystemModels(): any;
    accept(arg0?: any, arg1?: any): any;
    validateModelFile(arg0?: any, arg1?: any): any;
    _throwAlreadyExists(arg0?: any): any;
    addModelFile(arg0?: any, arg1?: any, arg2?: any): any;
    updateModelFile(arg0?: any, arg1?: any, arg2?: any): any;
    deleteModelFile(arg0?: any): any;
    addModelFiles(arg0?: any, arg1?: any, arg2?: any): any;
    validateModelFiles(): any;
    updateExternalModels(arg0?: any, arg1?: any): any;
    getModelFiles(): any;
    resolveType(arg0?: any, arg1?: any): any;
    clearModelFiles(): any;
    getModelFile(arg0?: any): any;
    getNamespaces(): any;
    getType(arg0?: any): any;
    getSystemTypes(): any;
    getAssetDeclarations(): any;
    getTransactionDeclarations(): any;
    getEventDeclarations(): any;
    getParticipantDeclarations(): any;
    getEnumDeclarations(): any;
    getConceptDeclarations(): any;
    getFactory(): any;
    getSerializer(): any;
  }
  export class ModelUtil {
    static getShortName(arg0?: any): any;
    static isWildcardName(arg0?: any): any;
    static isRecursiveWildcardName(arg0?: any): any;
    static isMatchingType(arg0?: any, arg1?: any): any;
    static getNamespace(arg0?: any): any;
    static getSystemNamespace(): any;
    static isPrimitiveType(arg0?: any): any;
    static isAssignableTo(arg0?: any, arg1?: any, arg2?: any): any;
    static capitalizeFirstLetter(arg0?: any): any;
    static isEnum(arg0?: any): any;
    static getFullyQualifiedName(arg0?: any, arg1?: any): any;
    constructor();
  }
  
  export class ParticipantDeclaration {
    constructor(arg0?: any, arg1?: any);
    isRelationshipTarget(): any;
    getSystemType(): any;
    validate(): any;
  }
  export class Property {
    constructor(arg0?: any, arg1?: any);
    getParent(): any;
    process(): any;
    validate(arg0?: any): any;
    getName(): any;
    getType(): any;
    isOptional(): any;
    getFullyQualifiedTypeName(): any;
    getFullyQualifiedName(): any;
    getNamespace(): any;
    isArray(): any;
    isTypeEnum(): any;
    isPrimitive(): any;
  }
  
  export class Relationship {
    static fromURI(arg0?: any, arg1?: any, arg2?: any, arg3?: any): any;
    constructor(arg0?: any, arg1?: any, arg2?: any, arg3?: any, arg4?: any);
    toString(): any;
    isRelationship(): any;
  }
  export class RelationshipDeclaration {
    constructor(arg0?: any, arg1?: any);
    validate(arg0?: any): any;
    toString(): any;
  }
  export class Resource {
    constructor(arg0?: any, arg1?: any, arg2?: any, arg3?: any, arg4?: any);
    toString(): any;
    isResource(): any;
    toJSON(): any;
  }
  
  export class TransactionDeclaration {
    constructor(arg0?: any, arg1?: any);
    getSystemType(): any;
    validate(): any;
  }
  
  
