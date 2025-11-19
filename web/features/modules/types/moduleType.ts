export interface ModuleVersion {
  version: string;
  date: string;
}

export interface TerraformModule {
  id: number;
  namespace: string;
  name: string;
  fullName: string;
  description: string;
  tags: string[];
  stars: number;
  version: string;
  provider: string;
  versions: ModuleVersion[];
}

export const modulesData: TerraformModule[] = [
  // ...your module objects here
];
