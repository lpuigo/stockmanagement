package adminmodal

const template string = `<el-dialog
		:visible.sync="visible" 
		:before-close="HideWithControl" :close-on-click-modal="false"
		width="80%" top="5vh"
>
	<!-- 	Modal Title	-->
	<span slot="title">
		<h2 style="margin: 0 0">
			<i class="fas fa-wrench icon--left"></i>Administration
		</h2>
	</span>

	<!-- 
		Modal Body
		style="height: 100%;"		
	-->
    <el-tabs type="border-card" tab-position="left" style="height: 70vh">
		<!-- ========================================== Admin Tab ================================================= -->
		<el-tab-pane label="Maintenance" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<el-button type="primary" @click="ReloadData" size="mini">Rechargement des données</el-button>
			<h4>&nbsp;</h4>
			<el-button type="primary" @click="SaveArchive" size="mini">Sauvegarde des archives</el-button>
<!--			<h4>&nbsp;</h4>-->
<!--			<h4>Archive des <a href="/api/worksites/archive">Chantiers Orange</a></h4>-->
		</el-tab-pane>

		<!-- ========================================== Users Tab ================================================= -->
		<el-tab-pane label="Utilisateurs" lazy=true style="height: 75vh; padding: 5px 25px; overflow-x: hidden;overflow-y: auto;">
			<el-table
					:border=false
					:data="filteredUsers"
					:default-sort = "{prop: 'Name', order: 'ascending'}"
					:row-class-name="TableRowClassName" height="90%" size="mini"
			>
				<!--	Edit User-->
				<el-table-column type="expand" width="40px" >
					<template slot-scope="scope">
						<!--	User Name & Password -->
						<el-row :gutter="5" type="flex" align="middle" class="spaced">
							<el-col :span="2" class="align-right">Nom:</el-col>
							<el-col :span="8">
								<el-input v-model="scope.row.Name" size="mini"></el-input>
							</el-col>
							<el-col :span="2" class="align-right">MdP:</el-col>
							<el-col :span="8">
								<el-input v-model="scope.row.Password" size="mini"></el-input>
							</el-col>
						</el-row>
						<!--	User Permission -->
						<el-row :gutter="5" type="flex" align="middle" class="spaced">
							<el-col :span="2" class="align-right">Permissions:</el-col>
							<el-col :span="6" >
								<p><el-switch v-model="scope.row.Permissions.Validate" active-text="Validation des mouvements"></el-switch></p>
							</el-col>
							<el-col :span="6" >
								<p><el-switch v-model="scope.row.Permissions.Invoice" active-text="Accès aux prix"></el-switch></p>
							</el-col>
							<el-col :span="6" >
								<p><el-switch v-model="scope.row.Permissions.Admin" active-text="Administrateur"></el-switch></p>
							</el-col>
						</el-row>
					</template>
				</el-table-column>

    			<!--	Index   -->
				<el-table-column
						label="N°" width="40px" align="right"
						type="index"
						index=1 
				></el-table-column>
			
				<!--	Actions   -->
				<el-table-column label="" width="80px">
					<template slot="header" slot-scope="scope">
						<el-tooltip  placement="top" open-delay=300 effect="light">
							<div slot="content">Ajouter un nouvel utilisateur</div>
							<el-button type="success" plain icon="fas fa-users fa-fw" size="mini" @click="AddNewUser()"></el-button>
						</el-tooltip>
					</template>
				</el-table-column>
				
				<!--	User Name   -->
				<el-table-column
						:resizable="true" :show-overflow-tooltip=true 
						prop="Name" label="Utilisateur" width="310px"
						sortable :sort-by="['Name']"
				></el-table-column>
				<!-- :filters="FilterList('Name')" :filter-method="FilterHandler"	filter-placement="bottom-end"-->
				
				<!--	Permissions   -->
				<el-table-column
						:resizable="true"
						label="Permissions"
				>
					<template slot-scope="scope">
						<span>
							<i class="fas fa-edit icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Validate }"></i>
							<i class="fas fa-euro-sign icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Invoice }"></i>
							<i class="fas fa-tools icon--medium icon--left" :class="{ 'icon--disabled': !scope.row.Permissions.Admin }"></i>
						</span>
					</template>
				</el-table-column>
			</el-table>
		</el-tab-pane>

    </el-tabs>

    <!-- 
        Modal Footer Action Bar
    -->
    <span slot="footer">
		<el-row :gutter="15">
			<el-col :span="24" style="text-align: right">
				<el-tooltip :open-delay="500" effect="light">
					<div slot="content">Annuler les changements</div>
					<el-button :disabled="!hasChanged" @click="UndoChange" icon="fas fa-undo-alt" plain size="mini"
                               type="info"></el-button>
				</el-tooltip>
				
				<el-button @click="Hide" size="mini">Fermer</el-button>
				
				<el-button :disabled="!hasChanged" type="success" @click="ConfirmChange" plain size="mini">Enregistrer</el-button>
			</el-col>
		</el-row>
	</span>
</el-dialog>`
